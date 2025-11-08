use bytes::Buf;
use futures::{Stream, TryStream, TryStreamExt, stream};
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
        mut stream: S, // The incoming byte stream (your 'field')
    ) -> Result<impl Stream<Item = Result<String, DocumentError>>, DocumentError>
    where
        S: TryStream<Ok = B, Error = E> + Unpin,
        B: Buf,
        E: Into<DocumentError>,
    {
        // --- 1. Consume the entire byte stream into a Vec<u8> ---
        let mut body = Vec::new();
        while let Some(chunk) = stream.try_next().await.map_err(Into::into)? {
            body.extend_from_slice(chunk.chunk());
        }

        // --- 2. Parse the PDF bytes from memory ---
        let extracted_text = match pdf_extract::extract_text_from_mem(&body) {
            Ok(text) => text,
            Err(e) => {
                tracing::error!("Failed to parse PDF: {}", e);
                // You might want a specific error here, like DocumentError::PdfParseFailed
                return Err(DocumentError::InternalError);
            }
        };

        // --- 3. Split the extracted text into lines ---
        // The embedding service expects a Vec<String> (lines)
        // TODO: curerntly we are just putting the entire text as a single part.
        // implement stream later
        let parts: Vec<String> = vec![extracted_text];

        // --- 4. Return a new stream from the Vec<String> ---
        // We wrap each String in a Result to match the required output signature
        let output_stream = stream::iter(parts.into_iter().map(Ok));

        Ok(output_stream)
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

    pub async fn get_document_by_id(
        &self,
        ctx: &RequestContext,
        document_id: Uuid,
    ) -> Result<Document, DocumentError> {
        let user_id = ctx.user_id.ok_or(DocumentError::InvalidCredentials)?;

        let document = self
            .document_repository
            .get_document_by_id_and_user_id(document_id, user_id)
            .await
            .map_err(|e| {
                tracing::error!("Error fetching document: {}", e);
                DocumentError::InternalError
            })?
            .ok_or(DocumentError::DocumentNotFound)?;

        Ok(document)
    }

    pub async fn delete_document_by_id(
        &self,
        ctx: &RequestContext,
        document_id: Uuid,
    ) -> Result<(), DocumentError> {
        let user_id = ctx.user_id.ok_or(DocumentError::InvalidCredentials)?;

        self.document_repository
            .delete_document_by_id_and_user_id(document_id, user_id)
            .await
            .map_err(|e| {
                tracing::error!("Error deleting document: {}", e);
                DocumentError::InternalError
            })?;

        Ok(())
    }
}
