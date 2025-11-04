use derive_more::derive::{Display, Error};
use uuid::Uuid;

#[derive(Debug, Display, Error)]
pub enum RefreshTokenPersistenceError {
    #[display("Unexpected error: {message}")]
    UnexpectedError { message: String },

    #[display("Token not found for user_id: {user_id}")]
    TokenNotFound { user_id: Uuid },
}
