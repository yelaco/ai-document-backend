use base64::{Engine as _, engine::general_purpose::STANDARD};
use chrono::{Duration, Utc};
use jsonwebtoken::{DecodingKey, EncodingKey, Header, encode};
use rand::RngCore;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use uuid::Uuid;

use crate::{
    application::auth::errors::AuthError,
    config::settings::Settings,
    core::{auth::Role, request_context::RequestContext},
    domain::{Auth, User},
    infrastructure::{auth::PasswordService, persistence::user::errors::UserPersistenceError},
    interfaces::UserRepository,
};

pub const ACCESS_TOKEN_EXPIRY_SECONDS: usize = 3600; // 1 hour

pub struct AuthService {
    repository: Arc<dyn UserRepository>,
    password_service: PasswordService,
    encoding_key: EncodingKey,
    decoding_key: DecodingKey,
}

impl AuthService {
    pub fn new(repository: Arc<dyn UserRepository>, settings: Settings) -> Self {
        Self {
            repository,
            password_service: PasswordService::new(),
            encoding_key: EncodingKey::from_secret(settings.auth.jwt_secret.as_ref()),
            decoding_key: DecodingKey::from_secret(settings.auth.jwt_secret.as_ref()),
        }
    }

    pub async fn register_user(
        &self,
        _ctx: RequestContext,
        email: &str,
        full_name: &str,
        password: &str,
    ) -> Result<User, AuthError> {
        let password_hash = self.password_service.hash_password(password).map_err(|e| {
            tracing::error!("Error hashing password: {}", e);
            AuthError::InternalError
        })?;

        self.repository
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
            .repository
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
        ctx: RequestContext,
        email: &str,
        password: &str,
    ) -> Result<Auth, AuthError> {
        tracing::info!(context = format!("{}", ctx));
        let Some(user) = self
            .repository
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

            // TODO: implement refresh token storage and management
            let refresh_token = self.create_refresh_token();

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
        STANDARD.encode(bytes)
    }

    pub fn validate_access_token(
        &self,
        token: &str,
    ) -> Result<Claims, jsonwebtoken::errors::Error> {
        let token_data = jsonwebtoken::decode::<Claims>(
            token,
            &self.decoding_key,
            &jsonwebtoken::Validation::default(),
        )?;
        Ok(token_data.claims)
    }
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Claims {
    pub sub: String,
    pub exp: usize,
    pub iss: String,
    pub role: String,
    pub email: String,
}
