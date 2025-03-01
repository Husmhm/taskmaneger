package user_test

import (
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

var bcryptGenerateFromPassword = bcrypt.GenerateFromPassword

func TestService_Register(t *testing.T) {
	mockValidator := new(MockValidator)
	mockRepo := new(MockRepository)
	mockAuth := new(MockAuth)

	s := user.New(mockRepo, mockValidator, mockAuth)
	validRequest := param.RegisterRequest{
		PhoneNumber: "+989131921277",
		Name:        "John Doe",
		Password:    "password",
	}

	t.Run("Invalid request", func(t *testing.T) {
		mockValidator.On("ValidateRegisterRequest", validRequest).Return(errors.New("validation error"))

		response, err := s.Register(validRequest)

		assert.Error(t, err)
		assert.Equal(t, param.RegisterResponse{}, response)
		mockValidator.AssertExpectations(t)
	})

	t.Run("Repository error", func(t *testing.T) {
		mockValidator.On("ValidateRegisterRequest", validRequest).Return(nil)

		expectedError := errors.New("repository error")
		mockRepo.On("Register", mock.AnythingOfType("models.User")).Return(models.User{}, expectedError)

		response, err := s.Register(validRequest)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected error")
		assert.Equal(t, param.RegisterResponse{}, response)
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Successful registration", func(t *testing.T) {
		mockValidator.On("ValidateRegisterRequest", validRequest).Return(nil)

		createdUser := models.User{
			Model:       gorm.Model{ID: 1},
			PhoneNumber: validRequest.PhoneNumber,
			Name:        validRequest.Name,
			Password:    "hashed_password",
		}
		mockRepo.On("Register", mock.AnythingOfType("models.User")).Return(createdUser, nil)

		response, err := s.Register(validRequest)

		assert.NoError(t, err)
		assert.Equal(t, param.RegisterResponse{
			User: param.UserInfo{
				ID:          createdUser.ID,
				PhoneNumber: createdUser.PhoneNumber,
				Name:        createdUser.Name,
			},
		}, response)
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
