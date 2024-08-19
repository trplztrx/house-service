package domain

import (
	"errors"
	"time"
)

type House struct {
	HouseID				int
	Address				string
	Year				int
	Developer			string
	CreateAt			time.Time
	UpdateAt			time.Time
}

var (
	ErrHouse_BadRequest = errors.New("bad house request for create")
	ErrHouse_BadID      = errors.New("bad house id")
	ErrHouse_BadYear    = errors.New("bad house construct year")
)