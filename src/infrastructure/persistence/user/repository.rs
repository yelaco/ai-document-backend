use crate::{
    domain::User,
    infrastructure::persistence::user::{errors::UserPersistenceError, row::UserRow},
    interfaces::UserRepository,
};
use async_trait::async_trait;
use sqlx::{Executor, Pool, Postgres, postgres::PgDatabaseError};
use uuid::Uuid;

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
    async fn create_user(
        &self,
        email: &str,
        password_hash: &str,
        full_name: &str,
    ) -> Result<u64, UserPersistenceError> {
        let result = self
            .pool
            .execute(sqlx::query!(
                r#"
            INSERT INTO users (email, password_hash, full_name)
            VALUES ($1, $2, $3)
                "#,
                email,
                password_hash,
                full_name
            ))
            .await
            .map_err(|e| match e {
                sqlx::Error::Database(db_err) => {
                    let pg_err = db_err.downcast_ref::<PgDatabaseError>();

                    if pg_err.code() == "23505" {
                        return UserPersistenceError::DuplicateUserError {
                            email: email.to_string(),
                        };
                    }

                    UserPersistenceError::UnexpectedError {
                        message: db_err.to_string(),
                    }
                }
                _ => UserPersistenceError::UnexpectedError {
                    message: e.to_string(),
                },
            })?;

        Ok(result.rows_affected())
    }

    async fn get_user_by_email(&self, email: &str) -> Result<Option<User>, String> {
        let user = sqlx::query_as!(
            UserRow,
            r#"
                SELECT id, email, password_hash, full_name, created_at, updated_at
            FROM users
            WHERE email = $1
            "#,
            email
        )
        .fetch_optional(&self.pool)
        .await
        .map(|opt_row| opt_row.map(User::from))
        .map_err(|e| e.to_string())?;

        Ok(user)
    }

    async fn get_user_by_id(&self, user_id: Uuid) -> Result<Option<User>, String> {
        let user = sqlx::query_as!(
            UserRow,
            r#"
            SELECT id, email, password_hash, full_name, created_at, updated_at
            FROM users
            WHERE id = $1
            "#,
            user_id
        )
        .fetch_optional(&self.pool)
        .await
        .map(|opt_row| opt_row.map(User::from))
        .map_err(|e| e.to_string())?;

        Ok(user)
    }
}
