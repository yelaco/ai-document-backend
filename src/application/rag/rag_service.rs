use qdrant_client::{
    Qdrant,
    qdrant::{
        Condition, Filter, PointId, QueryPointsBuilder, SearchParamsBuilder,
        point_id::PointIdOptions,
    },
};
use std::sync::Arc;
use uuid::Uuid;

use crate::application::rag::errors::RagError;

pub struct RagService {
    qdrant_client: Arc<Qdrant>,
}

impl RagService {
    pub fn new(qdrant_client: Arc<Qdrant>) -> Self {
        Self { qdrant_client }
    }

    pub async fn get_document_context(
        &self,
        document_id: Uuid,
        query_embedding: Vec<f32>,
        top_k: u64,
    ) -> Result<Vec<String>, RagError> {
        let document_chunk_ids = self
            .vector_search_document(document_id, query_embedding, top_k)
            .await
            .map_err(|e| {
                tracing::error!("Error getting document context: {}", e);
                RagError::InternalError
            })?;

        // TODO: Fetch the actual document chunks from persistence using the IDs

        Ok(vec!["hello world".to_string()])
    }

    async fn vector_search_document(
        &self,
        document_id: Uuid,
        query_embedding: Vec<f32>,
        top_k: u64,
    ) -> Result<Vec<Uuid>, RagError> {
        let search_result = self
            .qdrant_client
            .query(
                QueryPointsBuilder::new("documents")
                    .query(query_embedding)
                    .filter(Filter::must([Condition::matches(
                        "document_id",
                        document_id.to_string(),
                    )]))
                    .params(SearchParamsBuilder::default().hnsw_ef(128).exact(true))
                    .limit(top_k),
            )
            .await
            .map_err(|e| {
                tracing::error!("Error searching Qdrant: {}", e);
                RagError::InternalError
            })?;

        let ids: Vec<Uuid> = search_result
            .result
            .into_iter()
            .filter_map(|point| point.id?.point_id_options)
            .map(|point_id_options| match point_id_options {
                PointIdOptions::Num(id) => Uuid::from_u128(id as u128),
                PointIdOptions::Uuid(id) => Uuid::parse_str(&id).unwrap_or_else(|_| Uuid::nil()),
            })
            .collect();

        Ok(ids)
    }
}
