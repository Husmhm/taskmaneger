package user

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	models "taskmaneger/model"
	"taskmaneger/param"
	"taskmaneger/validator"
)

type Repository interface {
	Register(ctx context.Context, u models.User) (models.User, error)
}

type Service struct {
	repo      Repository
	validator validator.Validator
}

func New(repo Repository, validator validator.Validator) Service {
	return Service{repo: repo, validator: validator}
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
