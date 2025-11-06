pub struct PaginationParams {
    pub page: Option<u32>,
    pub per_page: Option<u32>,
}

#[derive(serde::Serialize)]
pub struct PaginationResponse<T> {
    pub items: Vec<T>,
    pub metadata: PaginationMetadata,
}

#[derive(serde::Serialize)]
pub struct PaginationMetadata {
    pub current_page: u32,
    pub per_page: u32,
}
