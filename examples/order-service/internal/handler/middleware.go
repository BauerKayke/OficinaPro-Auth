package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims representa os claims do JWT
type JWTClaims struct {
	UserID int64  `json:"userId"`
	Email  string `json:"email"`
	Nome   string `json:"nome"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ContextKey tipo para chaves de contexto
type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
	EmailKey  ContextKey = "email"
	RoleKey   ContextKey = "role"
	ClaimsKey ContextKey = "claims"
)

// JWTMiddleware valida JWT token e adiciona claims ao contexto
func JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrair token do header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			// Verificar formato "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondError(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validar token
			claims, err := validateJWT(tokenString, jwtSecret)
			if err != nil {
				respondError(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
				return
			}

			// Adicionar claims ao contexto
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			ctx = context.WithValue(ctx, ClaimsKey, claims)

			// Continuar com request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// validateJWT valida o token JWT e retorna os claims
func validateJWT(tokenString, jwtSecret string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extrair claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Verificar expiração
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, fmt.Errorf("token expired")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// RequireRole middleware para verificar role específica
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(RoleKey).(string)
			if !ok || userRole != role {
				respondError(w, "Insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID extrai o user ID do contexto
func GetUserID(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// GetUserEmail extrai o email do contexto
func GetUserEmail(ctx context.Context) (string, error) {
	email, ok := ctx.Value(EmailKey).(string)
	if !ok {
		return "", fmt.Errorf("email not found in context")
	}
	return email, nil
}

// GetUserRole extrai o role do contexto
func GetUserRole(ctx context.Context) (string, error) {
	role, ok := ctx.Value(RoleKey).(string)
	if !ok {
		return "", fmt.Errorf("role not found in context")
	}
	return role, nil
}

// respondError envia resposta de erro em JSON
func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":     message,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// respondJSON envia resposta em JSON
func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
