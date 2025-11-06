use chrono::{DateTime, Utc};
use sqlx::FromRow;
use uuid::Uuid;

use crate::{core::auth::Role, domain::User};

#[derive(Debug, FromRow)]
pub struct UserRow {
    pub id: Uuid,
    pub email: String,
    pub password_hash: String,
    pub full_name: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

impl From<UserRow> for User {
    fn from(row: UserRow) -> Self {
        User {
            id: row.id,
            email: row.email,
            password_hash: row.password_hash,
            full_name: row.full_name,
            created_at: row.created_at,
            updated_at: row.updated_at,
            role: Role::User,
        }
    }
}
