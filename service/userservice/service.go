package user

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	models "taskmaneger/model"
	"taskmaneger/param"
	"taskmaneger/service/authservice"
	"taskmaneger/validator"
)

type Repository interface {
	Register(ctx context.Context, u models.User) (models.User, error)
	GetUserByPhoneNumber(phone string) (models.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user models.User) (string, error)
	CreateRefreshToken(user models.User) (string, error)
}

type Service struct {
	repo      Repository
	validator validator.Validator
	auth      AuthGenerator
}

func New(repo Repository, validator validator.Validator, auth authservice.Service) Service {
	return Service{repo: repo, validator: validator, auth: auth}
}

func (s Service) Register(ctx context.Context, req param.RegisterRequest) (param.RegisterResponse, error) {
	vErr := s.validator.ValidateRegisterRequest(ctx, req)

	if vErr != nil {
		return param.RegisterResponse{}, vErr
	}

	password := []byte(req.Password)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	hashStr := string(hash)
	if err != nil {
		panic(err)
	}

	user := models.User{
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    hashStr,
	}
	createdUser, err := s.repo.Register(ctx, user)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return param.RegisterResponse{User: param.UserInfo{
		ID:          createdUser.ID,
		PhoneNumber: createdUser.PhoneNumber,
		Name:        createdUser.Name,
	}}, nil
}

func (s Service) Login(req param.LoginRequest) (param.LoginResponse, error) {
	// TODO - can use rich error
	vErr := s.validator.ValidateLoginRequest(req)

	if vErr != nil {
		return param.LoginResponse{}, vErr
	}

	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{}, err
	}

	hErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if hErr != nil {
		return param.LoginResponse{}, fmt.Errorf("invalid password")
	}
	// return ok
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return param.LoginResponse{User: param.UserInfo{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
	}, Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}}, nil

}
