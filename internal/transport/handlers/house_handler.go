package handlers

import (
	"context"
	"encoding/json"
	"house-service/internal/domain"
	"house-service/internal/transport/dto"
	"house-service/internal/usecase"
	"house-service/pkg"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)
type HouseHandler struct {
	usecase *usecase.HouseUsecase
	dbTimeout time.Duration
}

func NewHouseHandler(usecase *usecase.HouseUsecase, timeout time.Duration) *HouseHandler {
	return &HouseHandler{
		usecase: usecase,
		dbTimeout: timeout,
	}
}

//TODO: проверка запросов модерации и создания пользователем moderator
func (h *HouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		createRequest  transport.CreateHouseRequest
		createResponse transport.CreateHouseResponse
	)

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &createRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout * time.Second)
	defer cancel()

	createResponse, err = h.usecase.Create(ctx, &createRequest)
	if err != nil {
		http.Error(w, "Failed to create house", http.StatusInternalServerError)
		return
	}

	respBody, err := json.Marshal(createResponse)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

func (h *HouseHandler) GetApartmentsByID(w http.ResponseWriter, r *http.Request) {
	var getHouseResponse transport.GetHouseResponse

	houseID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout*time.Second)
	defer cancel()

	userRole, err := pkg.ExtractPayloadFromToken(r.Header.Get("authorization"), "role")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var status string
	if userRole == "moderator" {
		status = domain.AnyStatus
	} else {
		status = domain.ApprovedStatus
	}

	getHouseResponse, err = h.usecase.GetApartmentsByHouseID(ctx, houseID, status)
	if err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	respBody, err := json.Marshal(getHouseResponse)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

