package handler

import (
	"encoding/json"
	"net/http"

	"github.com/oficinapro/auth-service/internal/handler/response"
)

// ResponseBuilder facilita criação de HTTPResponses com JSON.
// Centraliza lógica de marshaling e headers.
type ResponseBuilder struct{}

// NewResponseBuilder cria um novo builder de responses.
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}

// Success cria response 200 OK com body JSON.
func (rb *ResponseBuilder) Success(data interface{}) HTTPResponse {
	return rb.buildResponse(http.StatusOK, data)
}

// BadRequest cria response 400 com mensagem de erro.
func (rb *ResponseBuilder) BadRequest(message string) HTTPResponse {
	return rb.buildErrorResponse(http.StatusBadRequest, message)
}

// Unauthorized cria response 401 com mensagem de erro.
func (rb *ResponseBuilder) Unauthorized(message string) HTTPResponse {
	return rb.buildErrorResponse(http.StatusUnauthorized, message)
}

// Forbidden cria response 403 com mensagem de erro.
func (rb *ResponseBuilder) Forbidden(message string) HTTPResponse {
	return rb.buildErrorResponse(http.StatusForbidden, message)
}

// NotFound cria response 404 com mensagem de erro.
func (rb *ResponseBuilder) NotFound(message string) HTTPResponse {
	return rb.buildErrorResponse(http.StatusNotFound, message)
}

// MethodNotAllowed cria response 405 com mensagem de erro.
func (rb *ResponseBuilder) MethodNotAllowed(message string) HTTPResponse {
	return rb.buildErrorResponse(http.StatusMethodNotAllowed, message)
}

// InternalError cria response 500 com mensagem genérica.
func (rb *ResponseBuilder) InternalError() HTTPResponse {
	return rb.buildErrorResponse(http.StatusInternalServerError, "Internal server error")
}

// CORSPreflight cria response para requests OPTIONS.
func (rb *ResponseBuilder) CORSPreflight() HTTPResponse {
	return HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       []byte{},
		Headers:    DefaultHeaders(),
	}
}

// buildResponse cria response com marshaling de data para JSON.
func (rb *ResponseBuilder) buildResponse(statusCode int, data interface{}) HTTPResponse {
	body, err := json.Marshal(data)
	if err != nil {
		// Fallback para erro interno se marshaling falhar
		return rb.InternalError()
	}

	return NewHTTPResponse(statusCode, body)
}

// buildErrorResponse cria response de erro com ErrorResponse padronizado.
func (rb *ResponseBuilder) buildErrorResponse(statusCode int, message string) HTTPResponse {
	errorResp := response.NewErrorResponse(message)
	body, err := json.Marshal(errorResp)
	if err != nil {
		// Fallback se marshaling falhar (raro para ErrorResponse simples)
		body = []byte(`{"error":"Internal server error","timestamp":""}`)
	}

	return NewHTTPResponse(statusCode, body)
}
