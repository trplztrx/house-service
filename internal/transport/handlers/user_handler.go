package handlers

import (
	"context"
	"encoding/json"
	"house-service/internal/transport/dto"
	"house-service/internal/usecase"
	"io"
	"net/http"
	"time"
)
type UserHandler struct {
	usecase *usecase.UserUsecase
	dbTimeout time.Duration
}

func NewUserHandler(usecase *usecase.UserUsecase, timeout time.Duration) *UserHandler {
	return &UserHandler{
		usecase: usecase,
		dbTimeout: timeout,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		registerRequest  transport.RegisterUserRequest
		registerResponse transport.RegisterUserResponse
	)

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &registerRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout * time.Second)
	defer cancel()

	registerResponse, err = h.usecase.Register(ctx, &registerRequest)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(registerResponse)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		loginRequest transport.LoginUserRequest
		loginResponse transport.LoginUserResponse
	)

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout * time.Second)
	defer cancel()

	loginResponse, err = h.usecase.Login(ctx, &loginRequest)
	if err != nil {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	respBody, err := json.Marshal(loginResponse)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

