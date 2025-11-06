use base64::{Engine as _, engine::general_purpose::URL_SAFE_NO_PAD};
use chrono::{Duration, Utc};
use jsonwebtoken::{DecodingKey, EncodingKey, Header, encode};
use rand::RngCore;
use std::sync::Arc;
use uuid::Uuid;

use crate::{
    application::auth::{claims::Claims, errors::AuthError},
    config::settings::Settings,
    core::{auth::Role, request_context::RequestContext},
    domain::{Auth, User},
    infrastructure::{
        auth::PasswordService,
        persistence::{
            refresh_token::errors::RefreshTokenPersistenceError, user::errors::UserPersistenceError,
        },
    },
    interfaces::{RefreshTokenRepository, UserRepository},
};

pub const ACCESS_TOKEN_EXPIRY_SECONDS: usize = 3600; // 1 hour

pub struct AuthService {
    user_repository: Arc<dyn UserRepository>,
    refresh_token_repository: Arc<dyn RefreshTokenRepository>,
    password_service: PasswordService,
    encoding_key: EncodingKey,
    decoding_key: DecodingKey,
}

impl AuthService {
    pub fn new(
        user_repository: Arc<dyn UserRepository>,
        refresh_token_repository: Arc<dyn RefreshTokenRepository>,
        settings: Settings,
    ) -> Self {
        Self {
            user_repository,
            refresh_token_repository,
            password_service: PasswordService::new(),
            encoding_key: EncodingKey::from_secret(settings.auth.jwt_secret.as_ref()),
            decoding_key: DecodingKey::from_secret(settings.auth.jwt_secret.as_ref()),
        }
    }

    pub async fn register_user(
        &self,
        _ctx: &RequestContext,
        email: &str,
        full_name: &str,
        password: &str,
    ) -> Result<User, AuthError> {
        let password_hash = self.password_service.hash_password(password).map_err(|e| {
            tracing::error!("Error hashing password: {}", e);
            AuthError::InternalError
        })?;

        self.user_repository
            .create_user(email, &password_hash, full_name)
            .await
            .map_err(|e| {
                tracing::error!("Error creating user: {}", e);
                match e {
                    UserPersistenceError::DuplicateUserError { email } => {
                        AuthError::UserAlreadyExists { email }
                    }
                    _ => AuthError::InternalError,
                }
            })?;

        let Some(user) = self
            .user_repository
            .get_user_by_email(email)
            .await
            .map_err(|e| {
                tracing::error!("Error fetching user by email: {}", e);
                AuthError::InternalError
            })?
        else {
            tracing::error!("New user not found");
            return Err(AuthError::InternalError);
        };

        Ok(user)
    }

    pub async fn login_user(
        &self,
        ctx: &RequestContext,
        email: &str,
        password: &str,
    ) -> Result<Auth, AuthError> {
        tracing::info!(context = format!("{}", ctx));
        let Some(user) = self
            .user_repository
            .get_user_by_email(email)
            .await
            .map_err(|e| {
                tracing::error!("Error fetching user by email: {}", e);
                AuthError::InternalError
            })?
        else {
            return Err(AuthError::UserNotFound);
        };

        if self
            .password_service
            .verify_password(password, &user.password_hash)
        {
            let access_token = self
                .create_access_token(&user.id, &user.email, &user.role)
                .map_err(|e| {
                    tracing::error!("Error creating access token: {}", e);
                    AuthError::InternalError
                })?;

            let refresh_token = self.create_refresh_token();
            let Ok(refresh_token_hash) = self.password_service.hash_password(&refresh_token) else {
                tracing::error!("Error hashing refresh token");
                return Err(AuthError::InternalError);
            };
            let new_expires_at = Utc::now()
                .checked_add_signed(Duration::days(30))
                .expect("time travel failed")
                .timestamp();

            self.refresh_token_repository
                .store_refresh_token(user.id, &refresh_token_hash, new_expires_at)
                .await
                .map_err(|e| {
                    tracing::error!("Error storing refresh token: {}", e);
                    AuthError::InternalError
                })?;

            Ok(Auth {
                access_token,
                refresh_token,
            })
        } else {
            Err(AuthError::InvalidCredentials)
        }
    }

    pub async fn refresh_flow(
        &self,
        ctx: &RequestContext,
        refresh_token: &str,
    ) -> Result<Auth, AuthError> {
        let user_id = ctx.user_id.ok_or(AuthError::InternalError)?;

        let token_hash = self
            .refresh_token_repository
            .get_refresh_token_hash(user_id)
            .await
            .map_err(|e| {
                tracing::error!("");
                match e {
                    RefreshTokenPersistenceError::TokenNotFound { user_id: _user_id } => {
                        AuthError::TokenValidationError
                    }
                    _ => AuthError::InternalError,
                }
            })?;

        if self
            .password_service
            .verify_password(refresh_token, &token_hash)
        {
            self.refresh_token_repository
                .revoke_refresh_token(user_id)
                .await
                .map_err(|e| {
                    tracing::error!("Error revoking refresh token: {}", e);
                    AuthError::InternalError
                })?;

            let Ok(Some(user)) = self.user_repository.get_user_by_id(user_id).await else {
                return Err(AuthError::UserNotFound);
            };

            let access_token = self
                .create_access_token(&user.id, &user.email, &user.role)
                .map_err(|e| {
                    tracing::error!("Error creating access token: {}", e);
                    AuthError::InternalError
                })?;

            let refresh_token = self.create_refresh_token();
            let Ok(refresh_token_hash) = self.password_service.hash_password(&refresh_token) else {
                tracing::error!("Error hashing refresh token");
                return Err(AuthError::InternalError);
            };
            let new_expires_at = Utc::now()
                .checked_add_signed(Duration::days(7))
                .expect("time travel failed")
                .timestamp();

            self.refresh_token_repository
                .store_refresh_token(user.id, &refresh_token_hash, new_expires_at)
                .await
                .map_err(|e| {
                    tracing::error!("Error storing refresh token: {}", e);
                    AuthError::InternalError
                })?;

            Ok(Auth {
                access_token,
                refresh_token,
            })
        } else {
            Err(AuthError::InvalidCredentials)
        }
    }

    fn create_access_token(
        &self,
        user_id: &Uuid,
        email: &str,
        role: &Role,
    ) -> Result<String, jsonwebtoken::errors::Error> {
        let expiration = Utc::now()
            .checked_add_signed(Duration::seconds(ACCESS_TOKEN_EXPIRY_SECONDS as i64))
            .expect("time travel failed")
            .timestamp() as usize;

        let claims = Claims {
            sub: user_id.to_string(),
            exp: expiration,
            iss: "ai-document-backend".to_string(),
            email: email.to_string(),
            role: role.to_string(),
        };

        let token = encode(&Header::default(), &claims, &self.encoding_key)?;

        Ok(token)
    }

    fn create_refresh_token(&self) -> String {
        let mut bytes = [0u8; 64];
        rand::rng().fill_bytes(&mut bytes);
        URL_SAFE_NO_PAD.encode(bytes)
    }

    pub fn validate_access_token(
        &self,
        token: &str,
        skip_exp_check: bool,
    ) -> Result<Claims, jsonwebtoken::errors::Error> {
        let mut validation = jsonwebtoken::Validation::default();
        if skip_exp_check {
            validation.validate_exp = false;
        }

        let token_data = jsonwebtoken::decode::<Claims>(token, &self.decoding_key, &validation)?;
        Ok(token_data.claims)
    }
}
