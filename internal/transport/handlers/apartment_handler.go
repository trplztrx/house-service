package handlers

import (
	"context"
	"encoding/json"
	transport "house-service/internal/transport/dto"
	"house-service/internal/usecase"
	"house-service/pkg"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ApartmentHandler struct {
	usecase  *usecase.ApartmentUsecase
	dbTimeout time.Duration
}

func NewApartmentHandler(usecase *usecase.ApartmentUsecase, timeout time.Duration) *ApartmentHandler {
	return &ApartmentHandler{
		usecase: usecase,
		dbTimeout: timeout,
	}
}

func (h *ApartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		apartmentRequest  transport.ApartmentCreateRequest
		apartmentResponse transport.ApartmentCreateResponse
	)

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &apartmentRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := pkg.ExtractPayloadFromToken(r.Header.Get("authorization"), "userID")
	if err != nil {
		http.Error(w, "Failed to extract user ID from token", http.StatusInternalServerError)
		return
	}
	userUuid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout*time.Second)
	defer cancel()

	apartmentResponse, err = h.usecase.Create(ctx, userUuid, &apartmentRequest)
	if err != nil {
		http.Error(w, "Failed to create flat", http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(apartmentResponse)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (h *ApartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		flatRequest  transport.ApartmentUpdateRequest
		flatResponse transport.ApartmentCreateResponse
	)

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &flatRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := pkg.ExtractPayloadFromToken(r.Header.Get("authorization"), "userID")
	if err != nil {
		http.Error(w, "Failed to extract user ID from token", http.StatusInternalServerError)
		return
	}
	userUuid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout*time.Second)
	defer cancel()

	flatResponse, err = h.usecase.Update(ctx, userUuid, &flatRequest)
	if err != nil {
		http.Error(w, "Failed to update flat", http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(flatResponse)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
