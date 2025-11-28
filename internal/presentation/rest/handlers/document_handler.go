package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/tasks"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/dtos"
)

type DocumentHandler struct {
	documentService interfaces.DocumentService
	taskDistributor tasks.TaskDistributor
}

func NewDocumentHandler(documentService interfaces.DocumentService, taskDistributor tasks.TaskDistributor) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
		taskDistributor: taskDistributor,
	}
}

func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(fmt.Errorf("invalid multipart form: %w", err))
		c.JSON(http.StatusBadRequest, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "invalid multipart form data",
			},
		})
		return
	}

	ext := filepath.Ext(file.Filename)
	tempFileName := fmt.Sprintf("%s_%s", uuid.New().String(), ext)
	savePath := filepath.Join("./tmp", tempFileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		_ = c.Error(fmt.Errorf("DocumentHandler.UploadDocument: failed to save uploaded file: %w", err))
		c.JSON(http.StatusInternalServerError, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "failed to save uploaded file",
			},
		})
		return
	}

	ctx := c.Request.Context()
	document, err := h.documentService.CreateDocument(ctx, file.Filename, savePath)
	if err != nil {
		_ = c.Error(fmt.Errorf("DocumentHandler.UploadDocument: failed to create document: %w", err))
		c.JSON(http.StatusInternalServerError, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "failed to create document",
			},
		})
		return
	}

	// distribute embedding task to processor
	taskPayload := tasks.PayloadEmbedDocument{
		DocumentID: document.ID,
		UserID:     document.UserID,
	}
	procOpt := tasks.TaskProcessingOption{
		Queue:     tasks.QueueDefault,
		MaxRetry:  5,
		ProcessIn: 10, // process after 10 seconds
	}
	_ = h.taskDistributor.DistributeTaskEmbedDocument(ctx, &taskPayload, procOpt)

	c.JSON(http.StatusOK, dtos.BaseAPIResponse{
		Status: dtos.StatusSuccess,
		Data:   dtos.DocumentResponseFromEntity(document),
	})
}

func (h *DocumentHandler) GetPaginatedDocuments(c *gin.Context) {
	var req dtos.PaginatedDocumentsParams
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(fmt.Errorf("DocumentHandler.GetPaginatedDocuments: failed to extract request body: %w", err))
		c.JSON(http.StatusBadRequest, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "invalid request payload",
			},
		})
		return
	}

	documents, count, err := h.documentService.GetPaginatedDocuments(
		c.Request.Context(),
		req.Page,
		req.PageSize,
	)
	if err != nil {
		_ = c.Error(fmt.Errorf("DocumentHandler.ListDocuments: failed to get documents: %w", err))
		c.JSON(http.StatusInternalServerError, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "failed to get documents",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseAPIResponse{
		Status: dtos.StatusSuccess,
		Data: dtos.PaginatedResponse[dtos.DocumentResponse]{
			Items: dtos.DocumentResponsesFromEntities(documents),
			Metadata: dtos.PaginationMetadata{
				TotalItems:  count,
				TotalPages:  (count + req.PageSize - 1) / req.PageSize,
				CurrentPage: req.Page,
				PageSize:    req.PageSize,
			},
		},
	})
}

func (h *DocumentHandler) GetDocumentByID(c *gin.Context) {
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(fmt.Errorf("DocumentHandler.GetDocumentByID: invalid document ID format: %w", err))
		c.JSON(http.StatusBadRequest, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "invalid document ID format",
			},
		})
		return
	}

	document, err := h.documentService.GetDocumentByID(c.Request.Context(), documentID)
	if err != nil {
		_ = c.Error(fmt.Errorf("DocumentHandler.GetDocumentByID: failed to get document: %w", err))
		c.JSON(http.StatusInternalServerError, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "failed to get document",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BaseAPIResponse{
		Status: dtos.StatusSuccess,
		Data:   dtos.DocumentResponseFromEntity(document),
	})
}

func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(fmt.Errorf("DocumentHandler.DeleteDocument: invalid document ID format: %w", err))
		c.JSON(http.StatusBadRequest, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "invalid document ID format",
			},
		})
		return
	}

	err = h.documentService.DeleteDocument(c.Request.Context(), documentID)
	if err != nil {
		_ = c.Error(fmt.Errorf("DocumentHandler.DeleteDocument: failed to delete document: %w", err))
		c.JSON(http.StatusInternalServerError, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "failed to delete document",
			},
		})
		return
	}

	// TODO: also delete document embeddings

	c.JSON(http.StatusOK, dtos.BaseAPIResponse{
		Status: dtos.StatusSuccess,
		Data:   gin.H{"message": "Document deleted successfully"},
	})
}
