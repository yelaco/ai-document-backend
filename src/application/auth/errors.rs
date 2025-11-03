use derive_more::derive::{Display, Error};

#[derive(Debug, Display, Error)]
pub enum AuthError {
    #[display("An internal error occurred. Please try again later.")]
    InternalError,

    #[display("User not found.")]
    UserNotFound,

    #[display("User already exists: {email}")]
    UserAlreadyExists { email: String },

    #[display("Invalid credentials provided.")]
    InvalidCredentials,

    #[display("Token creation failed")]
    TokenCreationError,

    #[display("Token validation failed")]
    TokenValidationError,
}
