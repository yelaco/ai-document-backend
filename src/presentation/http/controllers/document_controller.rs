use std::sync::Arc;

use actix_multipart::Multipart;
use actix_web::{HttpResponse, web};
use futures::TryStreamExt;
use sanitize_filename::sanitize;
use tokio::sync::Mutex;
use tracing::instrument;
use uuid::Uuid;

use crate::application::document::DocumentService;
use crate::application::document::errors::DocumentError;
use crate::application::embedding::EmbeddingService;
use crate::core::request_context::RequestContext;
use crate::presentation::http::dtos::{
    DocumentResponse, ListDocumentsRequest, PaginationMetadata, PaginationResponse,
};

#[tracing::instrument(name = "upload", skip(document_service, embedding_service, payload))]
pub async fn upload(
    ctx: RequestContext,
    document_service: web::Data<DocumentService>,
    embedding_service: web::Data<Arc<Mutex<EmbeddingService>>>,
    mut payload: Multipart,
) -> Result<HttpResponse, DocumentError> {
    let mut documents = Vec::new();

    while let Some(field) = payload.try_next().await.map_err(|_e| {
        tracing::error!("Error processing multipart payload");
        DocumentError::InternalError
    })? {
        let Some(content_disposition) = field.content_disposition() else {
            tracing::error!("Content disposition not found in multipart field");
            return Err(DocumentError::InternalError);
        };

        let filename = match content_disposition.get_filename() {
            Some(name) => sanitize(name),
            None => {
                tracing::error!("Filename not found in content disposition");
                return Err(DocumentError::InternalError);
            }
        };

        let document = document_service.create_document(&ctx, filename).await?;

        let document_content_stream = document_service.process_document(&ctx, field).await?;

        // TODO: use a background task for embedding
        embedding_service
            .lock()
            .await
            .embed_document(document.id, document_content_stream)
            .await
            .map_err(|e| {
                tracing::error!("Error embedding document: {}", e);
                DocumentError::InternalError
            })?;

        documents.push(document);
    }

    return Ok(HttpResponse::Ok().json(
        documents
            .iter()
            .map(DocumentResponse::from)
            .collect::<Vec<_>>(),
    ));
}

#[instrument(name = "get_paginated_documents", skip(document_service, query))]
pub async fn get_paginated_documents(
    ctx: RequestContext,
    document_service: web::Data<DocumentService>,
    query: web::Query<ListDocumentsRequest>,
) -> Result<HttpResponse, DocumentError> {
    let per_page = query.per_page.unwrap_or(10);
    let current_page = query.page.unwrap_or(1);

    let documents = document_service
        .get_paginated_documents(&ctx, current_page, per_page)
        .await
        .map_err(|e| {
            tracing::error!("Error fetching documents: {}", e);
            DocumentError::InternalError
        })?;

    Ok(HttpResponse::Ok().json(PaginationResponse::from(
        documents,
        PaginationMetadata {
            per_page,
            current_page,
        },
    )))
}

#[instrument(name = "get_document", skip(document_service))]
pub async fn get_document(
    ctx: RequestContext,
    document_service: web::Data<DocumentService>,
    document_id: web::Path<Uuid>,
) -> Result<HttpResponse, DocumentError> {
    let document = document_service
        .get_document_by_id(&ctx, document_id.into_inner())
        .await?;

    Ok(HttpResponse::Ok().json(DocumentResponse::from(document)))
}

#[instrument(name = "delete_document", skip(document_service, embedding_service))]
pub async fn delete_document(
    ctx: RequestContext,
    document_service: web::Data<DocumentService>,
    embedding_service: web::Data<Arc<Mutex<EmbeddingService>>>,
    document_id: web::Path<Uuid>,
) -> Result<HttpResponse, DocumentError> {
    let document_id = document_id.into_inner();
    document_service
        .delete_document_by_id(&ctx, document_id)
        .await?;

    let _ = embedding_service
        .lock()
        .await
        .delete_document_embeddings(document_id)
        .await
        .inspect_err(|e| {
            tracing::error!("Error deleting document embeddings: {}", e);
        });

    Ok(HttpResponse::NoContent().finish())
}
