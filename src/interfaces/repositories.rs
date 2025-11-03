use async_trait::async_trait;
use uuid::Uuid;

use crate::{domain::User, infrastructure::persistence::user::errors::UserPersistenceError};

#[async_trait]
pub trait UserRepository: Send + Sync {
    async fn create_user(
        &self,
        email: &str,
        password_hash: &str,
        full_name: &str,
    ) -> Result<u64, UserPersistenceError>;

    async fn get_user_by_email(&self, email: &str) -> Result<Option<User>, String>;

    async fn get_user_by_id(&self, user_id: Uuid) -> Result<Option<User>, String>;
}
