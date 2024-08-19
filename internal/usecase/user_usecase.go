package usecase

import (
	"context"
	"errors"
	"house-service/internal/domain"
	"house-service/internal/repo"
	"house-service/internal/transport/dto"
	"house-service/pkg"
	"github.com/google/uuid"
	
)

type UserUsecase struct {
	userRepo repo.UserRepo
}

func NewUserUsecase(userRepo repo.UserRepo) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) Register(ctx context.Context, registerRequest *transport.RegisterUserRequest) (response transport.RegisterUserResponse, err error) {
	// TODO: registerRequest validation
	existUser, err := u.userRepo.GetByEmail(ctx, registerRequest.Email)
	if err != nil {
		return
	}

	if existUser != nil {
		return response, errors.New("User already exist")
	}

	encryptedPassword, err := pkg.HashPassword(registerRequest.Password)
	if err != nil {
		return
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return
	}

	user := &domain.User{
		UserID: uuid,
		Email: registerRequest.Email,
		Password: encryptedPassword,
		Type: registerRequest.UserType,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return
	}

	token, err := pkg.GenerateToken(user.UserID, user.Type)  // TODO in pkg
	if err != nil {
		return
	}

	response = transport.RegisterUserResponse{Token: token}

	return
}

func (u *UserUsecase) Login(ctx context.Context, loginRequest *transport.LoginUserRequest) (response transport.LoginUserResponse, err error) {
	// TODO: loginRequest validation
	user, err := u.userRepo.GetByEmail(ctx, loginRequest.Email)
	if err != nil {
		return
	}

	if user == nil {
		return
	}

	err = pkg.CheckPassword(user.Password, loginRequest.Password)
	if err != nil {
		return
	}

	token, err := pkg.GenerateToken(user.UserID, user.Type)
	if err != nil {
		return
	}

	response = transport.LoginUserResponse{Token: token}

	return
}