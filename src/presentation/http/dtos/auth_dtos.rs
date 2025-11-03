use crate::domain::Auth;

#[derive(serde::Deserialize)]
pub struct LoginRequest {
    pub email: String,
    pub password: String,
}

#[derive(serde::Deserialize)]
pub struct RegisterRequest {
    pub email: String,
    pub full_name: String,
    pub password: String,
}

#[derive(serde::Serialize)]
pub struct LoginResponse {
    pub access_token: String,
    pub refresh_token: String,
}

impl LoginResponse {
    pub fn from(auth: Auth) -> Self {
        Self {
            access_token: auth.access_token,
            refresh_token: auth.refresh_token,
        }
    }
}
