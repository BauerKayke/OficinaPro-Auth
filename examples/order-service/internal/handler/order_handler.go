package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Order representa um pedido
type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	Status    string    `json:"status"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"createdAt"`
}

// HealthCheck retorna status do serviço
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{
		"status":  "ok",
		"service": "order-service",
		"time":    time.Now().UTC().Format(time.RFC3339),
	}, http.StatusOK)
}

// ListOrders lista pedidos do usuário autenticado
func ListOrders(w http.ResponseWriter, r *http.Request) {
	// Extrair user ID do contexto (vem do JWT middleware)
	userID, err := GetUserID(r.Context())
	if err != nil {
		respondError(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	// Simular busca no banco (substituir por repository real)
	orders := []Order{
		{
			ID:        1,
			UserID:    userID,
			Status:    "PENDING",
			Total:     150.50,
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        2,
			UserID:    userID,
			Status:    "COMPLETED",
			Total:     89.99,
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
	}

	respondJSON(w, orders, http.StatusOK)
}

// CreateOrder cria novo pedido
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Extrair user ID do contexto
	userID, err := GetUserID(r.Context())
	if err != nil {
		respondError(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req struct {
		Total float64 `json:"total"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validar
	if req.Total <= 0 {
		respondError(w, "Total must be greater than zero", http.StatusBadRequest)
		return
	}

	// Simular criação no banco
	order := Order{
		ID:        123, // Seria gerado pelo DB
		UserID:    userID,
		Status:    "PENDING",
		Total:     req.Total,
		CreatedAt: time.Now(),
	}

	respondJSON(w, order, http.StatusCreated)
}

// GetOrder obtém pedido específico
func GetOrder(w http.ResponseWriter, r *http.Request) {
	// Extrair user ID do contexto
	userID, err := GetUserID(r.Context())
	if err != nil {
		respondError(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	// Extrair order ID da URL
	vars := mux.Vars(r)
	orderID := vars["id"]

	// Simular busca no banco
	// Verificar se pedido pertence ao usuário
	order := Order{
		ID:        1,
		UserID:    userID,
		Status:    "PENDING",
		Total:     150.50,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}

	// Se order não pertence ao user, retornar 404
	// if order.UserID != userID { ... }

	respondJSON(w, order, http.StatusOK)
}
