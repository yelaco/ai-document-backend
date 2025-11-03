use std::sync::Arc;

use jsonwebtoken::{DecodingKey, EncodingKey};
use serde::{Deserialize, Serialize};

use crate::{
    application::auth::errors::AuthError,
    config::settings::Settings,
    core::request_context::RequestContext,
    domain::{Auth, User},
    infrastructure::auth::PasswordService,
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
                AuthError::InternalError
                // TODO: check for unique constraint violation
                // AuthError::UserAlreadyExists {
                //     email: email.to_string(),
                // }
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
            Ok(Auth {
                access_token: "mock_token".to_string(),
                refresh_token: "Bearer".to_string(),
            })
        } else {
            Err(AuthError::InvalidCredentials)
        }
    }
}

#[derive(Serialize, Deserialize)]
struct Claims {
    sub: String,
    exp: usize,
}
