use actix_multipart::MultipartError;
use derive_more::derive::{Display, Error};

#[derive(Debug, Display, Error)]
pub enum DocumentError {
    #[display("An internal error occurred. Please try again later.")]
    InternalError,

    #[display("Document not found.")]
    DocumentNotFound,

    #[display("Invalid credentials provided.")]
    InvalidCredentials,

    #[display("File upload error: {message}")]
    FileUploadError { message: String },
}

impl From<MultipartError> for DocumentError {
    fn from(err: MultipartError) -> Self {
        tracing::error!("Multipart error: {}", err);
        DocumentError::FileUploadError {
            message: err.to_string(),
        }
    }
}
