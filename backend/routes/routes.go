package routes

type ErrorResponse struct {
	ErrorMessages []string `json:"error_messages"`
}

func NewErrorResponse(messages []string) ErrorResponse {
	errorResponse := ErrorResponse{}
	errorResponse.ErrorMessages = messages
	return errorResponse
}
