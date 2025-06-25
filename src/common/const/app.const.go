package common
type Exception struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	HttpStatusCode int    `json:"httpStatusCode"`
}

// IServiceOutput<T> equivalent
type ServiceOutput[T any] struct {
	Message        string     `json:"message,omitempty"`
	OutputData     T          `json:"outputData,omitempty"`
	Exception      *Exception `json:"exception,omitempty"`
	HttpStatusCode int        `json:"httpStatusCode,omitempty"`
}

// Final API response structure
type ApiResponse[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}


const(Access_Token string = "access_token"
Refresh_Token string = "refresh_token")