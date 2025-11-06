use derive_more::derive::{Display, Error};

#[derive(Debug, Display, Error)]
pub enum DocumentPersistenceError {
    #[display("Unexpected error: {message}")]
    UnexpectedError { message: String },
}
