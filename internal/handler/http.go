package handler

// HTTPRequest abstrai request HTTP (framework agnostic)
type HTTPRequest struct {
	Method  string
	Body    []byte
	Headers map[string]string
}

// HTTPResponse abstrai response HTTP
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

// NewHTTPResponse cria response com defaults
func NewHTTPResponse(statusCode int, body []byte) HTTPResponse {
	return HTTPResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    DefaultHeaders(),
	}
}

// DefaultHeaders retorna headers padrão CORS
func DefaultHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
		"Access-Control-Allow-Methods": "POST,OPTIONS",
	}
}
