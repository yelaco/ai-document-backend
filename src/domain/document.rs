use chrono::{DateTime, Utc};
use uuid::Uuid;

use crate::domain::chat::Chat;

pub struct Document {
    pub id: Uuid,
    pub title: String,
    pub chats: Vec<Chat>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}
