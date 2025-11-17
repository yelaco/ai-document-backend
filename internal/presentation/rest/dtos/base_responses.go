package dtos

type BaseAPIResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type BaseErrorResponse struct {
	Status string         `json:"status"`
	Error  *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}
