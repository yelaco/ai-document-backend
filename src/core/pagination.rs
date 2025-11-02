pub struct PaginationParamsDto {
    pub page: u32,
    pub per_page: u32,
}

pub struct PaginationResponseDto<T> {
    pub items: Vec<T>,
    pub meta: PaginationMetaDto,
}

pub struct PaginationMetaDto {
    pub total_items: u32,
    pub total_pages: u32,
    pub current_page: u32,
    pub per_page: u32,
}
