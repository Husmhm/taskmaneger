package user_test

import (
	"github.com/stretchr/testify/mock"
	models "taskmaneger/model"
	"taskmaneger/param"
)

// Mock structures
type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) ValidateLoginRequest(req param.LoginRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockValidator) ValidateRegisterRequest(req param.RegisterRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByPhoneNumber(phoneNumber string) (models.User, error) {
	args := m.Called(phoneNumber)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockRepository) Register(u models.User) (models.User, error) {
	args := m.Called(u)
	return args.Get(0).(models.User), args.Error(1)
}

type MockAuth struct {
	mock.Mock
}

func (m *MockAuth) CreateAccessToken(user models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockAuth) CreateRefreshToken(user models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}
