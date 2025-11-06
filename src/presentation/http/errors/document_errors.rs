use crate::application::document::errors::DocumentError;
use actix_web::http::StatusCode;
use actix_web::http::header::ContentType;
use actix_web::{HttpResponse, ResponseError};

impl ResponseError for DocumentError {
    fn status_code(&self) -> StatusCode {
        match self {
            DocumentError::InternalError => StatusCode::INTERNAL_SERVER_ERROR,
            DocumentError::DocumentNotFound => StatusCode::NOT_FOUND,
            DocumentError::InvalidCredentials => StatusCode::UNAUTHORIZED,
            DocumentError::FileUploadError { message: _message } => StatusCode::BAD_REQUEST,
        }
    }

    fn error_response(&self) -> HttpResponse {
        HttpResponse::build(self.status_code())
            .insert_header(ContentType::html())
            .body(self.to_string())
    }
}
