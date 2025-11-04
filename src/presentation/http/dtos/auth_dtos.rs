use crate::domain::Auth;

#[derive(serde::Deserialize, Debug)]
pub struct LoginRequest {
    pub email: String,
    pub password: String,
}

#[derive(serde::Deserialize, Debug)]
pub struct RegisterRequest {
    pub email: String,
    pub full_name: String,
    pub password: String,
}

#[derive(serde::Serialize)]
pub struct TokenResponse {
    pub access_token: String,
}

impl TokenResponse {
    pub fn from(auth: Auth) -> Self {
        Self {
            access_token: auth.access_token,
        }
    }
}
