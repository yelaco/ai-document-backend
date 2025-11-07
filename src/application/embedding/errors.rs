use derive_more::derive::{Display, Error};

#[derive(Debug, Display, Error)]
pub enum EmbeddingError {
    #[display("An internal error occurred. Please try again later.")]
    InternalError,

    #[display("Invalid document ID provided.")]
    InvalidDocumentId,
}
