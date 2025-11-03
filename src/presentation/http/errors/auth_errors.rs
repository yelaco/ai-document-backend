use crate::application::auth::errors::AuthError;
use actix_web::http::StatusCode;
use actix_web::http::header::ContentType;
use actix_web::{HttpResponse, ResponseError};

impl ResponseError for AuthError {
    fn status_code(&self) -> StatusCode {
        match self {
            AuthError::InternalError => StatusCode::INTERNAL_SERVER_ERROR,
            AuthError::UserNotFound => StatusCode::NOT_FOUND,
            AuthError::InvalidCredentials => StatusCode::UNAUTHORIZED,
            AuthError::TokenCreationError => StatusCode::INTERNAL_SERVER_ERROR,
            AuthError::TokenValidationError => StatusCode::UNAUTHORIZED,
            AuthError::UserAlreadyExists { email: _email } => StatusCode::CONFLICT,
        }
    }

    fn error_response(&self) -> HttpResponse {
        HttpResponse::build(self.status_code())
            .insert_header(ContentType::html())
            .body(self.to_string())
    }
}
