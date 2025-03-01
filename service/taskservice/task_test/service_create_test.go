package task_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	models "taskmaneger/model"
	"taskmaneger/param"
	"taskmaneger/service/taskservice"
	"testing"
	"time"
)

func TestService_Create(t *testing.T) {

	userID := uint(1)

	t.Run("Successful task creation", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		validRequest := param.CreateTaskRequest{
			Title:       "Test Task",
			Description: "This is a test task",
		}

		createdTask := models.Task{
			ID:          1,
			Title:       validRequest.Title,
			Description: validRequest.Description,
			Status:      0,
			UserID:      userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("CreateTask", mock.AnythingOfType("models.Task")).Return(createdTask, nil)

		response, err := service.Create(validRequest, userID)

		assert.NoError(t, err)
		assert.Equal(t, param.CreateTaskResponse{
			Id:     createdTask.ID,
			Title:  createdTask.Title,
			Status: createdTask.Status,
			UserID: createdTask.UserID,
		}, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Task creation error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		validRequest := param.CreateTaskRequest{
			Title:       "Test Task",
			Description: "This is a test task",
		}

		expectedError := errors.New("repository error")
		mockRepo.On("CreateTask", mock.AnythingOfType("models.Task")).Return(models.Task{}, expectedError)

		response, err := service.Create(validRequest, userID)

		assert.Error(t, err)
		assert.Equal(t, param.CreateTaskResponse{}, response)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
