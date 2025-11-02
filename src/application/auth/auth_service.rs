use std::sync::Arc;

use jsonwebtoken::{DecodingKey, EncodingKey};
use serde::{Deserialize, Serialize};

use crate::{
    config::settings::Settings,
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

    pub fn register_user(
        &self,
        _email: &str,
        _full_name: &str,
        _password: &str,
    ) -> Result<User, String> {
        // Registration logic to be implemented
        unimplemented!()
    }

    pub async fn login_user(&self, _email: &str, _password: &str) -> Result<Auth, String> {
        Ok(Auth {
            access_token: "mock_token".to_string(),
            refresh_token: "Bearer".to_string(),
        })
    }

    pub async fn verify_basic_credentials(
        &self,
        email: &str,
        password: &str,
    ) -> Result<Option<User>, String> {
        let user = self.repository.get_user_by_email(email).await?;

        let Some(user) = user else {
            return Ok(None);
        };

        if self
            .password_service
            .verify_password(password, &user.password_hash)
        {
            Ok(Some(user))
        } else {
            Ok(None)
        }
    }
}

#[derive(Serialize, Deserialize)]
struct Claims {
    sub: String,
    exp: usize,
}
