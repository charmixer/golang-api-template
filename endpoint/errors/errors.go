package errors

/*
Codes we care about

200 OK
400 BadRequest
401 Unauthorized - Authentication denied
403 Forbidden - Authorization denined
404 Not Found
500 Internal Server Error
503 Service Unavailable
*/

type FieldError struct {
	Path string `json:"path"`
	Err string `json:"err"`
}

type HttpClientErrorResponse struct {
	StatusCode int `json:"status_code"`
	Method string `json:"method"`
	Url string `json:"url"`
	Errors []FieldError `json:"errors"`
}
