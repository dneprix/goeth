package proto

type SendEthRequest struct {
	//Name      string `json:"name"`
}

type SendEthResponse struct {
	//Name      string `json:"name"`
}

type GetLastResponse struct {
	//Name      string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
	}
}
