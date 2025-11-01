use chrono::{DateTime, Utc};
use uuid::Uuid;

pub struct User {
    pub id: Uuid,
    pub email: String,
    pub full_name: String,
    pub password_hash: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

impl User {
    pub fn new(email: String, full_name: String, password_hash: String) -> Self {
        Self {
            id: Uuid::new_v4(),
            email,
            full_name,
            password_hash,
            created_at: Utc::now(),
            updated_at: Utc::now(),
        }
    }
}
