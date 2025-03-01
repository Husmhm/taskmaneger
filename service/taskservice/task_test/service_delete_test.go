package task_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	models "taskmaneger/model"
	"taskmaneger/param"
	"taskmaneger/service/taskservice"
	"testing"
)

func TestService_Delete(t *testing.T) {
	taskID := uint(1)
	userID := uint(1)
	validRequest := param.DeleteTaskRequest{
		TaskId: taskID,
	}

	t.Run("Successful task deletion", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		existingTask := models.Task{
			ID:     taskID,
			Title:  "Test Task",
			UserID: userID,
		}

		mockRepo.On("GetTaskById", taskID).Return(existingTask, nil)
		mockRepo.On("DeleteTask", taskID).Return(nil)

		response, err := service.Delete(validRequest, userID)

		assert.NoError(t, err)
		assert.Equal(t, param.DeleteTaskResponse{Message: "successfully delete"}, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Task not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		expectedError := errors.New("task not found")
		mockRepo.On("GetTaskById", taskID).Return(models.Task{}, expectedError)

		response, err := service.Delete(validRequest, userID)

		assert.Error(t, err)
		assert.Equal(t, param.DeleteTaskResponse{}, response)
		assert.EqualError(t, err, expectedError.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("User not authorized", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		existingTask := models.Task{
			ID:     taskID,
			Title:  "Test Task",
			UserID: userID + 1,
		}

		mockRepo.On("GetTaskById", taskID).Return(existingTask, nil)

		response, err := service.Delete(validRequest, userID)

		assert.Error(t, err)
		assert.Equal(t, param.DeleteTaskResponse{}, response)
		assert.EqualError(t, err, "not authorized")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete task error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		existingTask := models.Task{
			ID:     taskID,
			Title:  "Test Task",
			UserID: userID,
		}

		expectedError := errors.New("delete error")
		mockRepo.On("GetTaskById", taskID).Return(existingTask, nil)
		mockRepo.On("DeleteTask", taskID).Return(expectedError)

		response, err := service.Delete(validRequest, userID)

		assert.Error(t, err)
		assert.Equal(t, param.DeleteTaskResponse{}, response)
		assert.EqualError(t, err, expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}
