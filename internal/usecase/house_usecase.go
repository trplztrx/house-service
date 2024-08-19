package usecase

import (
	"context"
	"errors"
	"time"
	"house-service/internal/domain"
	"house-service/internal/repo"
	"house-service/internal/transport/dto"
)
type HouseUsecase struct {
	houseRepo repo.HouseRepo
}

func NewHouseUsecase(houseRepo repo.HouseRepo) *HouseUsecase {
	return &HouseUsecase{
		houseRepo: houseRepo,
	}
}

func (u *HouseUsecase) Create(ctx context.Context, createRequest *transport.CreateHouseRequest) (response transport.CreateHouseResponse, err error) {
	//TODO: houseRequest validation

	existHouse, err := u.houseRepo.GetByAddress(ctx, createRequest.Address)
	if err != nil {
		return
	}

	if existHouse != nil {
		return response, errors.New("House already exist")
	}

	house := &domain.House{
		HouseID: createRequest.HouseID,
		Address: createRequest.Address,
		Year: createRequest.Year,
		Developer: createRequest.Developer,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	_, err = u.houseRepo.Create(ctx, house)
	if err != nil {
		return
	}

	response = transport.CreateHouseResponse{
		HouseID: createRequest.HouseID,
		Address: createRequest.Address,
		Year: createRequest.Year,
		Developer: createRequest.Developer,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	return
}

func (u *HouseUsecase) GetApartmentsByHouseID(ctx context.Context, houseID int, status string) (response transport.GetHouseResponse, err error) {

	flats, err := u.houseRepo.GetApartmentsByID(ctx, houseID, status)
	if err != nil {
		return
	}

	response = transport.GetHouseResponse{
		Apartments: flats,
	}

	return
}


