use futures::TryStream;
use futures_util::TryStreamExt;
use std::sync::Arc;
use uuid::Uuid;

use crate::{
    application::document::errors::DocumentError, core::request_context::RequestContext,
    domain::Document, interfaces::DocumentRepository,
};

pub struct DocumentService {
    document_repository: Arc<dyn DocumentRepository>,
}

impl DocumentService {
    pub fn new(document_repository: Arc<dyn DocumentRepository>) -> Self {
        Self {
            document_repository,
        }
    }

    pub async fn process_document<S, B, E>(
        &self,
        _ctx: &RequestContext,
        mut stream: S,
        document_id: Uuid,
    ) -> Result<(), DocumentError>
    where
        S: TryStream<Ok = B, Error = E> + Unpin,
        B: AsRef<[u8]>,
        E: Into<DocumentError>,
    {
        println!("Processing document with ID: {}", document_id);

        while let Some(chunk) = stream.try_next().await.map_err(Into::into)? {
            let data = chunk.as_ref();
            // Process the chunk (e.g., save to storage, analyze content, etc.)
            println!("Processing chunk of size: {}", data.len());
        }

        Ok(())
    }

    pub async fn create_document(
        &self,
        ctx: &RequestContext,
        title: String,
    ) -> Result<Document, DocumentError> {
        let user_id = ctx.user_id.ok_or(DocumentError::InvalidCredentials)?;

        let document = self
            .document_repository
            .create_document(title, user_id)
            .await
            .map_err(|e| {
                tracing::error!("Error creating document: {}", e);
                DocumentError::InternalError
            })?;

        Ok(document)
    }

    pub async fn get_paginated_documents(
        &self,
        ctx: &RequestContext,
        page: u32,
        page_size: u32,
    ) -> Result<Vec<Document>, DocumentError> {
        let user_id = ctx.user_id.ok_or(DocumentError::InvalidCredentials)?;

        let documents = self
            .document_repository
            .get_paginated_documents_by_user_id(user_id, page, page_size)
            .await
            .map_err(|e| {
                tracing::error!("Error fetching documents: {}", e);
                DocumentError::InternalError
            })?;

        Ok(documents)
    }
}
