package dtos

type PaginationParams struct {
	Page     int64 `form:"page" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=1,max=100"`
}

type PaginationMetadata struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int64 `json:"total_pages"`
	CurrentPage int64 `json:"current_page"`
	PageSize    int64 `json:"page_size"`
}

type PaginatedResponse[T any] struct {
	Items    []T                `json:"items"`
	Metadata PaginationMetadata `json:"metadata"`
}
