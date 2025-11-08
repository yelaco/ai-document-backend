use derive_more::derive::{Display, Error};

#[derive(Debug, Display, Error)]
pub enum RagError {
    #[display("An internal error occurred. Please try again later.")]
    InternalError,
}
