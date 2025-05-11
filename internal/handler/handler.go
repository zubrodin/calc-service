package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zubrodin/calc-service/internal/auth"
	"github.com/zubrodin/calc-service/internal/repository"
	"github.com/zubrodin/calc-service/internal/service"
)

type Handler struct {
	service *service.Service
	repo    repository.Repository
}

func New(s *service.Service, repo repository.Repository) *Handler {
	return &Handler{
		service: s,
		repo:    repo,
	}
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	_, err := h.repo.CreateUser(req.Login, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrUserExists {
			status = http.StatusConflict
		}
		respondWithError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	user, err := h.repo.Authenticate(req.Login, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrUserNotFound || err == repository.ErrInvalidPassword {
			status = http.StatusUnauthorized
		}
		respondWithError(w, status, "Invalid credentials")
		return
	}

	token, err := auth.GenerateToken(user.Login)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{Token: token})
}
func (h *Handler) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing token")
			return
		}

		_, err := auth.ValidateToken(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Invalid request format")
		return
	}

	result, err := h.service.Calculate(req.Expression)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidExpression {
			status = http.StatusUnprocessableEntity
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, CalculateResponse{Result: result})
}
