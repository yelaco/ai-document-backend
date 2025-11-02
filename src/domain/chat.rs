use chrono::{DateTime, Utc};
use uuid::Uuid;

pub struct Chat {
    pub id: Uuid,
    pub title: Option<String>,
    pub document_id: Uuid,
    pub user_id: Uuid,
    pub messages: Vec<Uuid>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}
