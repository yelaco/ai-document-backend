use chrono::{DateTime, Utc};
use uuid::Uuid;

pub struct Message {
    pub id: Uuid,
    pub role: String,
    pub content: String,
    pub chat_id: Uuid,
    pub timestamp: DateTime<Utc>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}
