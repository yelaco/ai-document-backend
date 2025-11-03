use derive_more::derive::{Display, Error};

#[derive(Debug, Display, Error)]
pub enum UserPersistenceError {
    #[display("Duplicate user error for {email}")]
    DuplicateUserError { email: String },

    #[display("Unexpected error: {message}")]
    UnexpectedError { message: String },
}
