package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Handler wraps auth service
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes sets up auth routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Get("/me", h.Me)
	})
}

// Register request struct
type registerRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	AccountID string `json:"account_id"`
	Role      string `json:"role"`
}

// Login request struct
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handler
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.Register(r.Context(), RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		AccountID: req.AccountID,
		Role:      req.Role,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": token})
}

// Login handler
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": token})
}

// Me handler
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate JWT and get claims
	claims, err := h.service.jwt.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Extract ID from claims and pass to service
	user, err := h.service.Me(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return basic user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":         user.ID,
		"email":      user.Email,
		"account_id": user.AccountID,
		"role":       user.Role,
	})
}
