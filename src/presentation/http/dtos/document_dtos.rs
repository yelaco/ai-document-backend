use serde::{Deserialize, Serialize};

use crate::{
    domain::Document,
    presentation::http::dtos::pagination::{PaginationMetadata, PaginationResponse},
};

#[derive(Deserialize)]
pub struct ListDocumentsRequest {
    pub page: Option<u32>,
    pub per_page: Option<u32>,
}

#[derive(Serialize)]
pub struct DocumentResponse {
    pub id: String,
    pub title: String,
    pub user_id: String,
    pub created_at: String,
    pub updated_at: String,
}

impl From<&Document> for DocumentResponse {
    fn from(doc: &Document) -> Self {
        Self {
            id: doc.id.to_string(),
            title: doc.title.clone(),
            user_id: doc.user_id.to_string(),
            created_at: doc.created_at.to_rfc3339(),
            updated_at: doc.updated_at.to_rfc3339(),
        }
    }
}

impl From<Document> for DocumentResponse {
    fn from(doc: Document) -> Self {
        Self {
            id: doc.id.to_string(),
            title: doc.title,
            user_id: doc.user_id.to_string(),
            created_at: doc.created_at.to_rfc3339(),
            updated_at: doc.updated_at.to_rfc3339(),
        }
    }
}

impl PaginationResponse<DocumentResponse> {
    pub fn from(documents: Vec<Document>, metadata: PaginationMetadata) -> Self {
        let items = documents
            .into_iter()
            .map(|doc| DocumentResponse {
                id: doc.id.to_string(),
                title: doc.title,
                user_id: doc.user_id.to_string(),
                created_at: doc.created_at.to_rfc3339(),
                updated_at: doc.updated_at.to_rfc3339(),
            })
            .collect();

        Self { items, metadata }
    }
}
