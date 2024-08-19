package usecase

import (
	"context"
	"errors"
	"house-service/internal/domain"
	"house-service/internal/repo"
	"house-service/internal/transport/dto"

	"github.com/google/uuid"
)

type ApartmentUsecase struct {
	apartmentRepo repo.ApartmentRepo
	houseRepo repo.HouseRepo
}

func NewApartmentUsecase(apartmentRepo repo.ApartmentRepo, houseRepo repo.HouseRepo) *ApartmentUsecase {
	return &ApartmentUsecase{
		apartmentRepo: apartmentRepo,
		houseRepo: houseRepo,
	}
}

func (u *ApartmentUsecase) Create(ctx context.Context, userID uuid.UUID, createRequest *transport.ApartmentCreateRequest) (response transport.ApartmentCreateResponse, err error) {
	// TODO: createRequest validation

	// TODO: exist flat in db

	house, err := u.houseRepo.GetByID(ctx, createRequest.HouseID)
	if err != nil {
		return
	}

	if house == nil {
		return response, errors.New("House does not exist")
	}


	apartment := &domain.Apartment{
		ApartmentID: createRequest.ApartmentID,
		HouseID: createRequest.HouseID,
		OwnerID: userID,
		Price:	createRequest.Price,
		Rooms: createRequest.Rooms,
		Status: "created",
	}

	err = u.houseRepo.Update(ctx, house)
	if err != nil {
		return
	}

	response = transport.ApartmentCreateResponse{
		ApartmentID: apartment.ApartmentID,
		HouseID: apartment.HouseID,
		Price: apartment.Price,
		Rooms: apartment.Rooms,
		Status: apartment.Status,
	}

	return
}

func (u *ApartmentUsecase) Update(ctx context.Context, moderatorID uuid.UUID, updateRequest *transport.ApartmentUpdateRequest) (response transport.ApartmentCreateResponse, err error) {

	apartment := &domain.Apartment {
		ApartmentID: updateRequest.ApartmentID,
		HouseID: updateRequest.HouseID,
		Status:  updateRequest.Status,
	}

	updatedApartment, err := u.apartmentRepo.Update(ctx, moderatorID, apartment)
	if err != nil {
		return
	}

	response = transport.ApartmentCreateResponse{
		ApartmentID: updatedApartment.ApartmentID,
		HouseID: updatedApartment.HouseID,
		Price: updatedApartment.Price,
		Rooms: updatedApartment.Rooms,
		Status: updatedApartment.Status,
	}

	return
}