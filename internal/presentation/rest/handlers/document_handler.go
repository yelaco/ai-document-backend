package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
)

type DocumentHandler struct {
	documentService interfaces.DocumentService
}

func NewDocumentHandler(documentService interfaces.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

func (h *DocumentHandler) UploadDocument(c *gin.Context) {
}

func (h *DocumentHandler) ListDocuments(c *gin.Context) {
}

func (h *DocumentHandler) GetDocumentByID(c *gin.Context) {}

func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
}
