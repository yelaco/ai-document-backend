use chrono::{DateTime, Utc};
use sqlx::FromRow;
use uuid::Uuid;

use crate::domain::Document;

#[derive(Debug, FromRow)]
pub struct DocumentRow {
    pub id: Uuid,
    pub title: String,
    pub user_id: Uuid,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

impl From<DocumentRow> for Document {
    fn from(row: DocumentRow) -> Self {
        Document {
            id: row.id,
            title: row.title,
            user_id: row.user_id,
            created_at: row.created_at,
            updated_at: row.updated_at,
        }
    }
}
