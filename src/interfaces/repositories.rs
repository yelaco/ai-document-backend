use async_trait::async_trait;
use uuid::Uuid;

use crate::domain::Document;
use crate::infrastructure::persistence::document::errors::DocumentPersistenceError;
use crate::infrastructure::persistence::refresh_token::errors::RefreshTokenPersistenceError;
use crate::{domain::User, infrastructure::persistence::user::errors::UserPersistenceError};

#[async_trait]
pub trait UserRepository: Send + Sync {
    async fn create_user(
        &self,
        email: &str,
        password_hash: &str,
        full_name: &str,
    ) -> Result<u64, UserPersistenceError>;

    async fn get_user_by_email(&self, email: &str) -> Result<Option<User>, String>;

    async fn get_user_by_id(&self, user_id: Uuid) -> Result<Option<User>, String>;
}

#[async_trait]
pub trait RefreshTokenRepository: Send + Sync {
    async fn store_refresh_token(
        &self,
        user_id: Uuid,
        refresh_token: &str,
        expires_at: i64,
    ) -> Result<u64, RefreshTokenPersistenceError>;

    async fn get_refresh_token_hash(
        &self,
        user_id: Uuid,
    ) -> Result<String, RefreshTokenPersistenceError>;

    async fn revoke_refresh_token(&self, user_id: Uuid)
    -> Result<(), RefreshTokenPersistenceError>;

    async fn delete_refresh_token(&self, user_id: Uuid)
    -> Result<(), RefreshTokenPersistenceError>;
}

#[async_trait]
pub trait DocumentRepository: Send + Sync {
    async fn create_document(
        &self,
        title: String,
        user_id: Uuid,
    ) -> Result<Document, DocumentPersistenceError>;

    async fn get_document_by_id(
        &self,
        document_id: Uuid,
    ) -> Result<Option<Document>, DocumentPersistenceError>;

    async fn get_paginated_documents_by_user_id(
        &self,
        user_id: Uuid,
        page: u32,
        page_size: u32,
    ) -> Result<Vec<Document>, DocumentPersistenceError>;
}
