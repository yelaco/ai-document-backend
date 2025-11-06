use chrono::{DateTime, Utc};
use uuid::Uuid;

use crate::core::auth::Role;

pub struct User {
    pub id: Uuid,
    pub email: String,
    pub full_name: String,
    pub password_hash: String,
    pub role: Role,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}
