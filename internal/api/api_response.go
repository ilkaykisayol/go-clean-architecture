package api

type ApiResponse struct {
	Data    *interface{} `json:"Data"`
	Code    byte         `json:"Code"`
	Message string       `json:"Message"`
}

func Ok(data interface{}) *ApiResponse {
	apiResponse := ApiResponse{
		Data:    &data,
		Code:    0,
		Message: "Success",
	}

	return &apiResponse
}

func Error(message string) *ApiResponse {
	apiResponse := ApiResponse{
		Data:    nil,
		Code:    1,
		Message: message,
	}

	return &apiResponse
}
