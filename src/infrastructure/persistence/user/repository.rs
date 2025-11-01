use crate::{domain::User, interfaces::UserRepository};
use async_trait::async_trait;
use sqlx::{Pool, Postgres};

pub struct PostgresUserRepository {
    pool: Pool<Postgres>,
}

impl PostgresUserRepository {
    pub fn new(pool: Pool<Postgres>) -> Self {
        Self { pool }
    }
}

#[async_trait]
impl UserRepository for PostgresUserRepository {
    fn create_user(
        &self,
        email: &str,
        password_hash: &str,
        full_name: &str,
    ) -> Result<u32, String> {
        Ok(1)
    }

    fn get_user_by_email(&self, email: &str) -> Result<Option<User>, String> {
        Ok(None)
    }

    fn get_user_by_id(&self, user_id: &str) -> Result<Option<User>, String> {
        Ok(None)
    }
}
