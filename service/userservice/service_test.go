package user_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	models "taskmaneger/model"
	"taskmaneger/param"
	user "taskmaneger/service/userservice"
	"testing"
)

// Mock structures
type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) ValidateLoginRequest(req param.LoginRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockValidator) ValidateRegisterRequest(ctx context.Context, req param.RegisterRequest) error {
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

func (m *MockRepository) Register(ctx context.Context, u models.User) (models.User, error) {
	args := m.Called(ctx, u)
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

func TestService_Login(t *testing.T) {
	mockValidator := new(MockValidator)
	mockRepo := new(MockRepository)
	mockAuth := new(MockAuth)

	s := user.New(mockRepo, mockValidator, mockAuth)

	t.Run("Invalid request", func(t *testing.T) {
		req := param.LoginRequest{}
		mockValidator.On("ValidateLoginRequest", req).Return(errors.New("invalid request"))

		response, err := s.Login(req)

		assert.Error(t, err)
		assert.Equal(t, param.LoginResponse{}, response)
		mockValidator.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		req := param.LoginRequest{PhoneNumber: "989131921299", Password: "password"}
		mockValidator.On("ValidateLoginRequest", req).Return(nil)
		mockRepo.On("GetUserByPhoneNumber", req.PhoneNumber).Return(models.User{}, errors.New("user not found"))

		response, err := s.Login(req)

		assert.Error(t, err)
		assert.Equal(t, param.LoginResponse{}, response)
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid password", func(t *testing.T) {
		req := param.LoginRequest{PhoneNumber: "989131921277", Password: "wrong_password"}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("hash_pass"), bcrypt.DefaultCost)
		if err != nil {
			t.Fatal(err)
		}

		user := models.User{PhoneNumber: "989131921277", Password: string(hashedPassword)}
		mockValidator.On("ValidateLoginRequest", req).Return(nil)
		mockRepo.On("GetUserByPhoneNumber", req.PhoneNumber).Return(user, nil)

		response, err := s.Login(req)

		assert.Error(t, err)
		assert.Equal(t, param.LoginResponse{}, response)
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Token creation error", func(t *testing.T) {
		req := param.LoginRequest{PhoneNumber: "+989131921277", Password: "password"}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			t.Fatal(err)
		}

		user := models.User{
			Model: gorm.Model{
				ID: 1,
			},
			PhoneNumber: "989131921277", Password: string(hashedPassword),
		}

		mockValidator.On("ValidateLoginRequest", req).Return(nil)
		mockRepo.On("GetUserByPhoneNumber", req.PhoneNumber).Return(user, nil)
		mockAuth.On("CreateAccessToken", user).Return("", errors.New("token creation error"))

		response, err := s.Login(req)

		assert.Error(t, err)
		assert.Equal(t, param.LoginResponse{}, response)
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("Successful login", func(t *testing.T) {
		req := param.LoginRequest{PhoneNumber: "1234567890", Password: "password"}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			t.Fatal(err)
		}

		user := models.User{
			Model: gorm.Model{
				ID: 1,
			},
			PhoneNumber: "989131921277", Password: string(hashedPassword), Name: "John Doe",
		}

		mockValidator.On("ValidateLoginRequest", req).Return(nil)
		mockRepo.On("GetUserByPhoneNumber", req.PhoneNumber).Return(user, nil)
		mockAuth.On("CreateAccessToken", user).Return("access_token", nil)
		mockAuth.On("CreateRefreshToken", user).Return("refresh_token", nil)

		response, err := s.Login(req)

		assert.NoError(t, err)
		assert.Equal(t, param.LoginResponse{
			User: param.UserInfo{
				ID:          1,
				PhoneNumber: "989131921277",
				Name:        "John Doe",
			},
			Tokens: param.Tokens{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
			},
		}, response)
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})
}
