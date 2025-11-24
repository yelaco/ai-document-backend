package dtos

import "github.com/gin-gonic/gin"

type ResponseStatus string

const (
	StatusSuccess ResponseStatus = "success"
	StatusError   ResponseStatus = "error"
)

type BaseAPIResponse struct {
	Status ResponseStatus `json:"status"`
	Data   gin.H          `json:"data"`
}

type BaseErrorResponse struct {
	Status ResponseStatus `json:"status"`
	Error  ErrorResponse  `json:"error"`
}

type ErrorResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}
