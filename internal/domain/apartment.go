package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Apartment struct {
	ApartmentID int
	HouseID 	int
	OwnerID 	uuid.UUID
	ModeratorID int
	Price 		int
	Rooms 		int
	Status 		string
}

const (
	CreatedStatus    = "created"
	DeclinedStatus   = "declined"
	ApprovedStatus   = "approved"
	ModeratingStatus = "on moderation"
	AnyStatus        = "any"
)

var (
	ErrFlat_BadPrice   = errors.New("bad flat price")
	ErrFlat_BadID      = errors.New("bad flat id")
	ErrFlat_BadHouseID = errors.New("bad flats house id")
	ErrFlat_BadRooms   = errors.New("bad flat rooms")
	ErrFlat_BadNewFlat = errors.New("bad new flat for update")
	ErrFlat_BadStatus  = errors.New("bad flat status")
	ErrFlat_BadRequest = errors.New("bad request for create")
)