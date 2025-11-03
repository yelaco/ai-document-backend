use async_trait::async_trait;

use crate::domain::User;

#[async_trait]
pub trait UserRepository: Send + Sync {
    async fn create_user(
        &self,
        email: &str,
        password_hash: &str,
        full_name: &str,
    ) -> Result<u64, String>;

    async fn get_user_by_email(&self, email: &str) -> Result<Option<User>, String>;

    async fn get_user_by_id(&self, user_id: &str) -> Result<Option<User>, String>;
}
