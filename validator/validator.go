package validator

import "context"

const (
	phoneNumberRegex = `^\+989[0-9]{9}$`
)

type Repository interface {
	IsPhoneNumberUnique(ctx context.Context, phoneNumber string) (bool, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}
