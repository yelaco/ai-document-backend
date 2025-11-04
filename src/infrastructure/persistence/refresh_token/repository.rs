use async_trait::async_trait;
use sqlx::{Pool, Postgres, postgres::PgDatabaseError};
use uuid::Uuid;

use crate::{
    infrastructure::persistence::refresh_token::errors::RefreshTokenPersistenceError,
    interfaces::RefreshTokenRepository,
};

pub struct PostgresRefreshTokenRepository {
    pool: Pool<Postgres>,
}

impl PostgresRefreshTokenRepository {
    pub fn new(pool: Pool<Postgres>) -> Self {
        Self { pool }
    }
}

#[async_trait]
impl RefreshTokenRepository for PostgresRefreshTokenRepository {
    async fn store_refresh_token(
        &self,
        user_id: Uuid,
        refresh_token_hash: &str,
        expires_at: i64,
    ) -> Result<u64, RefreshTokenPersistenceError> {
        let result = sqlx::query!(
            r#"
            INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
            VALUES ($1, $2, to_timestamp($3))
            "#,
            user_id,
            refresh_token_hash,
            expires_at as f64
        )
        .execute(&self.pool)
        .await
        .map_err(|e| match e {
            sqlx::Error::Database(db_err) => {
                let pg_err = db_err.downcast_ref::<PgDatabaseError>();

                if pg_err.code() == "23505" {
                    unimplemented!();
                }

                RefreshTokenPersistenceError::UnexpectedError {
                    message: db_err.to_string(),
                }
            }
            _ => RefreshTokenPersistenceError::UnexpectedError {
                message: e.to_string(),
            },
        })?;

        Ok(result.rows_affected())
    }

    async fn get_refresh_token_hash(
        &self,
        user_id: Uuid,
    ) -> Result<String, RefreshTokenPersistenceError> {
        let Some(record) = sqlx::query!(
            r#"
            SELECT id, token_hash FROM refresh_tokens
            WHERE user_id = $1  AND revoked = false AND expires_at > NOW()
            ORDER BY updated_at DESC
            "#,
            user_id,
        )
        .fetch_optional(&self.pool)
        .await
        .map_err(|e| RefreshTokenPersistenceError::UnexpectedError {
            message: e.to_string(),
        })?
        else {
            return Err(RefreshTokenPersistenceError::TokenNotFound { user_id });
        };

        Ok(record.token_hash)
    }

    async fn revoke_refresh_token(
        &self,
        user_id: Uuid,
    ) -> Result<(), RefreshTokenPersistenceError> {
        sqlx::query!(
            r#"
            UPDATE refresh_tokens
            SET revoked = true, revoked_at = NOW()
            WHERE user_id = $1 AND revoked = false
            "#,
            user_id,
        )
        .execute(&self.pool)
        .await
        .map_err(|e| RefreshTokenPersistenceError::UnexpectedError {
            message: e.to_string(),
        })?;

        Ok(())
    }

    async fn delete_refresh_token(
        &self,
        user_id: Uuid,
    ) -> Result<(), RefreshTokenPersistenceError> {
        sqlx::query!(
            r#"
            DELETE FROM refresh_tokens
            WHERE user_id = $1
            "#,
            user_id,
        )
        .execute(&self.pool)
        .await
        .map_err(|e| RefreshTokenPersistenceError::UnexpectedError {
            message: e.to_string(),
        })?;

        Ok(())
    }
}
