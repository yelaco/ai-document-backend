use std::io;

use actix_multipart::MultipartError;
use derive_more::derive::{Display, Error};
use tokio_util::codec::LinesCodecError;

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

    #[display("IO error: {message}")]
    Io { message: String },

    #[display("The uploaded document has a line that exceeds the maximum allowed length.")]
    LineTooLong,
}

impl From<MultipartError> for DocumentError {
    fn from(err: MultipartError) -> Self {
        tracing::error!("Multipart error: {}", err);
        DocumentError::FileUploadError {
            message: err.to_string(),
        }
    }
}

impl From<io::Error> for DocumentError {
    fn from(e: io::Error) -> Self {
        DocumentError::Io {
            message: e.to_string(),
        }
    }
}

impl From<DocumentError> for std::io::Error {
    fn from(e: DocumentError) -> Self {
        // We convert our custom error into a generic I/O error.
        // .to_string() ensures the error message is preserved.
        std::io::Error::other(e)
    }
}

impl From<LinesCodecError> for DocumentError {
    fn from(e: LinesCodecError) -> Self {
        match e {
            // This re-uses your existing `From<io::Error>` logic
            LinesCodecError::Io(io_err) => io_err.into(),

            // This maps the new error case to your new variant
            LinesCodecError::MaxLineLengthExceeded => DocumentError::LineTooLong,
        }
    }
}
