package transport

import (
	"house-service/internal/domain"
	"time"
)

type CreateHouseRequest struct {
	HouseID 	int 	`json:"id"`
	Address   	string 	`json:"address"`
	Year		int    	`json:"year"`
	Developer 	string 	`json:"developer"`
}

type CreateHouseResponse struct {
	HouseID 	int 		`json:"id"`
	Address   	string 		`json:"address"`
	Year		int    		`json:"year"`
	Developer 	string 		`json:"developer"`
	CreateAt	time.Time	`json:"created_at"`
	UpdateAt	time.Time	`json:"updated_at"`
}

type GetHouseResponse struct {
    // HouseID        	string          `json:"house_id"`
    // Address        	string          `json:"address"`
    // Year     		int             `json:"year"`
    // Developer      	string          `json:"developer,omitempty"` // Поле может быть пустым, поэтому используем `omitempty`
    Apartments 		[]domain.Apartment  `json:"apartments"`
    // CreateAt 		string          `json:"created_at"`
    // UpdateAt 		string          `json:"last_added_at"`
}