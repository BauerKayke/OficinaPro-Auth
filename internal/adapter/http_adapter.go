package adapter

import (
	"context"
	"io"
	"net/http"

	"github.com/oficinapro/auth-service/internal/handler"
)

// HTTPHandler interface para handler HTTP (permite mocking).
type HTTPHandler interface {
	Handle(ctx context.Context, req handler.HTTPRequest) handler.HTTPResponse
}

// HTTPAdapter adapta http.Handler padrão para handler.HTTPRequest/Response.
// Implementa o padrão Adapter para desacoplar framework HTTP do handler.
type HTTPAdapter struct {
	handler HTTPHandler
}

// NewHTTPAdapter cria um novo adapter HTTP.
func NewHTTPAdapter(h HTTPHandler) *HTTPAdapter {
	return &HTTPAdapter{
		handler: h,
	}
}

// ServeHTTP implementa http.Handler para compatibilidade com net/http.
func (a *HTTPAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Converte http.Request → handler.HTTPRequest
	httpReq := a.convertRequest(r)

	// Processa com handler (framework agnostic)
	httpResp := a.handler.Handle(r.Context(), httpReq)

	// Escreve resposta HTTP
	a.writeResponse(w, httpResp)
}

// convertRequest converte http.Request para handler.HTTPRequest.
func (a *HTTPAdapter) convertRequest(r *http.Request) handler.HTTPRequest {
	// Lê body
	var body []byte
	if r.Body != nil {
		defer r.Body.Close()
		body, _ = io.ReadAll(r.Body)
	}

	// Converte headers (pega primeiro valor de cada header)
	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return handler.HTTPRequest{
		Method:  r.Method,
		Path:    r.URL.Path,
		Body:    body,
		Headers: headers,
	}
}

// writeResponse escreve handler.HTTPResponse para http.ResponseWriter.
func (a *HTTPAdapter) writeResponse(w http.ResponseWriter, resp handler.HTTPResponse) {
	// Define headers
	for k, v := range resp.Headers {
		w.Header().Set(k, v)
	}

	// Define status code
	w.WriteHeader(resp.StatusCode)

	// Escreve body
	if resp.Body != nil {
		w.Write(resp.Body)
	}
}
