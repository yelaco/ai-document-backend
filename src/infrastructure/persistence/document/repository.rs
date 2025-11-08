use async_trait::async_trait;
use sqlx::{Pool, Postgres};
use uuid::Uuid;

use crate::{
    domain::Document,
    infrastructure::persistence::document::{errors::DocumentPersistenceError, row::DocumentRow},
    interfaces::DocumentRepository,
};

pub struct PostgresDocumentRepository {
    pool: Pool<Postgres>,
}

impl PostgresDocumentRepository {
    pub fn new(pool: Pool<Postgres>) -> Self {
        Self { pool }
    }
}

#[async_trait]
impl DocumentRepository for PostgresDocumentRepository {
    async fn create_document(
        &self,
        title: String,
        user_id: Uuid,
    ) -> Result<Document, DocumentPersistenceError> {
        let document_id = sqlx::query_scalar!(
            r#"
                INSERT INTO documents (title, user_id)
                VALUES ($1, $2)
                RETURNING id
                "#,
            title,
            user_id
        )
        .fetch_one(&self.pool)
        .await
        .map_err(|e| DocumentPersistenceError::UnexpectedError {
            message: e.to_string(),
        })?;

        self.get_document_by_id_and_user_id(document_id, user_id)
            .await?
            .ok_or_else(|| DocumentPersistenceError::UnexpectedError {
                message: "Failed to retrieve the created document".to_string(),
            })
    }

    async fn get_document_by_id_and_user_id(
        &self,
        document_id: Uuid,
        user_id: Uuid,
    ) -> Result<Option<Document>, DocumentPersistenceError> {
        let row = sqlx::query_as!(
            DocumentRow,
            r#"
            SELECT id, title, user_id, created_at, updated_at
            FROM documents
            WHERE id = $1 AND user_id = $2
            "#,
            document_id,
            user_id
        )
        .fetch_optional(&self.pool)
        .await
        .map_err(|e| DocumentPersistenceError::UnexpectedError {
            message: e.to_string(),
        })?;

        Ok(row.map(Document::from))
    }

    async fn get_paginated_documents_by_user_id(
        &self,
        user_id: Uuid,
        page: u32,
        page_size: u32,
    ) -> Result<Vec<Document>, DocumentPersistenceError> {
        let rows = sqlx::query_as!(
            DocumentRow,
            r#"
            SELECT id, title, user_id, created_at, updated_at
            FROM documents
            WHERE user_id = $1
            ORDER BY created_at DESC
            LIMIT $2 OFFSET $3
            "#,
            user_id,
            page_size as i64,
            ((page - 1) * page_size) as i64
        )
        .fetch_all(&self.pool)
        .await
        .map_err(|e| DocumentPersistenceError::UnexpectedError {
            message: e.to_string(),
        })?;

        Ok(rows.into_iter().map(Document::from).collect())
    }

    async fn delete_document_by_id_and_user_id(
        &self,
        document_id: Uuid,
        user_id: Uuid,
    ) -> Result<(), DocumentPersistenceError> {
        sqlx::query!(
            r#"
            DELETE FROM documents
            WHERE id = $1 AND user_id = $2
            "#,
            document_id,
            user_id
        )
        .execute(&self.pool)
        .await
        .map_err(|e| DocumentPersistenceError::UnexpectedError {
            message: e.to_string(),
        })?;

        Ok(())
    }
}
