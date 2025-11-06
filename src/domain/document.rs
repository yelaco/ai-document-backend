use chrono::{DateTime, Utc};
use uuid::Uuid;

pub struct Document {
    pub id: Uuid,
    pub title: String,
    pub user_id: Uuid,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}
