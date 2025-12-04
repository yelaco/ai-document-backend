package dtos

type CreateDocumentParams struct {
	OriginalName string
	SavePath     string
}

type UpdateDocumentParams struct {
	Title       *string
	Description *string
	FileURL     *string
}
