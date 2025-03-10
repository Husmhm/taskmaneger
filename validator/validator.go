package validator

import (
	models "taskmaneger/model"
)

const (
	phoneNumberRegex = `^\+989[0-9]{9}$`
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (models.User, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}
