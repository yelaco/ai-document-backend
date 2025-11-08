use std::sync::Arc;

use fastembed::{EmbeddingModel, InitOptions, TextEmbedding};
use futures::{Stream, StreamExt};
use qdrant_client::{
    Payload, Qdrant,
    qdrant::{PointStruct, UpsertPointsBuilder, Vector, Vectors},
};
use text_splitter::{ChunkConfig, TextSplitter};
use uuid::Uuid;

use crate::application::{document::errors::DocumentError, embedding::errors::EmbeddingError};

pub struct EmbeddingService {
    qdrant_client: Arc<Qdrant>,
}

impl EmbeddingService {
    pub fn new(qdrant_client: Arc<Qdrant>) -> Self {
        Self { qdrant_client }
    }

    pub async fn embed_texts(&self, texts: Vec<String>) -> Result<Vec<Vec<f32>>, EmbeddingError> {
        let mut embedder = TextEmbedding::try_new(
            InitOptions::new(EmbeddingModel::MultilingualE5Base)
                .with_show_download_progress(true)
                .with_max_length(768),
        )
        .expect("Failed to create embedder");

        let embeddings = embedder.embed(texts, None).map_err(|e| {
            tracing::error!("Error generating embeddings: {}", e);
            EmbeddingError::InternalError
        })?;

        Ok(embeddings)
    }

    pub async fn embed_document<S>(
        &self,
        document_id: Uuid,
        mut document_chunks: S,
    ) -> Result<(), EmbeddingError>
    where
        S: Stream<Item = Result<String, DocumentError>> + Unpin,
    {
        let splitter = TextSplitter::new(ChunkConfig::new(500).with_overlap(100).map_err(|e| {
            tracing::error!("Error creating text splitter: {}", e);
            EmbeddingError::InternalError
        })?);

        let mut embedder = TextEmbedding::try_new(
            InitOptions::new(EmbeddingModel::MultilingualE5Base)
                .with_show_download_progress(true)
                .with_max_length(768),
        )
        .expect("Failed to create embedder");

        while let Some(c) = document_chunks.next().await {
            match c {
                Ok(document_chunk) => {
                    if document_chunk.trim().is_empty() {
                        continue;
                    }
                    let content_chunks = splitter.chunks(&document_chunk);

                    let embeddings =
                        embedder
                            .embed(content_chunks.collect(), None)
                            .map_err(|e| {
                                tracing::error!("Error generating embeddings: {}", e);
                                EmbeddingError::InternalError
                            })?;

                    self.qdrant_client
                        .upsert_points(
                            UpsertPointsBuilder::new(
                                "documents",
                                embeddings
                                    .into_iter()
                                    .map(|vector| PointStruct {
                                        // TODO: store chunks in database and put their IDs here
                                        id: Some(uuid::Uuid::new_v4().to_string().into()),
                                        vectors: Some(Vectors::from(Vector::from(vector))),
                                        payload: Payload::try_from(serde_json::json!({
                                            "document_id": document_id,
                                        }))
                                        .unwrap()
                                        .into(),
                                    })
                                    .collect::<Vec<PointStruct>>(),
                            )
                            .wait(true),
                        )
                        .await
                        .map_err(|e| {
                            tracing::error!("Error embedding document: {}", e);
                            EmbeddingError::InternalError
                        })?;
                }
                Err(e) => {
                    tracing::error!("Error reading document chunk: {}", e);
                }
            }
        }

        Ok(())
    }

    pub async fn delete_document_embeddings(
        &self,
        _document_id: Uuid,
    ) -> Result<(), EmbeddingError> {
        unimplemented!()
    }
}
