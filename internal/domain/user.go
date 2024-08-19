package domain

import (
	"errors"

	"github.com/google/uuid"
)

const (
	Client string = "client"
	Moderator string = "moderator"
)

type User struct {
	UserID 		uuid.UUID
	Email 		string
	Password 	string
	Type		string
}

var (
	ErrUser_BadType     = errors.New("bd user type")
	ErrUser_BadRequest  = errors.New("bad nil request")
	ErrUser_BadMail     = errors.New("bad mail")
	ErrUser_BadPassword = errors.New("bad password")
)

