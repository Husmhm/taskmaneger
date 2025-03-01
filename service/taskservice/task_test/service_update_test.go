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

func TestService_Update(t *testing.T) {
	taskID := uint(1)
	userID := uint(1)

	t.Run("Successful task update", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		validRequest := param.UpdateTaskRequest{
			TaskId:      taskID,
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      models.StatusInprogress,
		}

		existingTask := models.Task{
			ID:     taskID,
			Title:  "Test Task",
			UserID: userID,
		}

		updatedTask := models.Task{
			ID:          taskID,
			Title:       validRequest.Title,
			Description: validRequest.Description,
			Status:      validRequest.Status,
			UserID:      userID,
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("GetTaskById", taskID).Return(existingTask, nil)
		mockRepo.On("UpdateTask", mock.MatchedBy(func(task models.Task) bool {
			return task.ID == taskID &&
				task.Title == validRequest.Title &&
				task.Description == validRequest.Description &&
				task.Status == validRequest.Status &&
				task.UserID == userID
		})).Return(updatedTask, nil)

		response, err := service.Update(validRequest, userID)

		assert.NoError(t, err)
		assert.Equal(t, param.UpdateTaskResponse{Task: updatedTask}, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Task not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		validRequest := param.UpdateTaskRequest{
			TaskId:      taskID,
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      models.StatusInprogress,
		}

		expectedError := errors.New("task not found")
		mockRepo.On("GetTaskById", taskID).Return(models.Task{}, expectedError)

		response, err := service.Update(validRequest, userID)

		assert.Error(t, err)
		assert.Equal(t, param.UpdateTaskResponse{}, response)
		assert.EqualError(t, err, expectedError.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("User not authorized", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		validRequest := param.UpdateTaskRequest{
			TaskId:      taskID,
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      models.StatusInprogress,
		}

		existingTask := models.Task{
			ID:     taskID,
			Title:  "Test Task",
			UserID: userID + 1,
		}

		mockRepo.On("GetTaskById", taskID).Return(existingTask, nil)

		response, err := service.Update(validRequest, userID)

		assert.Error(t, err)
		assert.Equal(t, param.UpdateTaskResponse{}, response)
		assert.EqualError(t, err, "not authorized")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid status", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		existingTask := models.Task{
			ID:     taskID,
			Title:  "Test Task",
			UserID: userID,
		}

		invalidRequest := param.UpdateTaskRequest{
			TaskId:      taskID,
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      4,
		}

		mockRepo.On("GetTaskById", taskID).Return(existingTask, nil)

		response, err := service.Update(invalidRequest, userID)

		assert.Error(t, err)
		assert.Equal(t, param.UpdateTaskResponse{}, response)
		assert.EqualError(t, err, "invalid status")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update task error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		validRequest := param.UpdateTaskRequest{
			TaskId:      taskID,
			Title:       "Updated Task",
			Description: "Updated Description",
			Status:      models.StatusInprogress,
		}

		existingTask := models.Task{
			ID:     taskID,
			Title:  "Test Task",
			UserID: userID,
			Status: models.StatusTodo,
		}

		expectedError := errors.New("update error")
		mockRepo.On("GetTaskById", taskID).Return(existingTask, nil)
		mockRepo.On("UpdateTask", mock.MatchedBy(func(task models.Task) bool {
			return task.ID == taskID &&
				task.Title == validRequest.Title &&
				task.Description == validRequest.Description &&
				task.Status == validRequest.Status &&
				task.UserID == userID
		})).Return(models.Task{}, expectedError)

		response, err := service.Update(validRequest, userID)

		assert.Error(t, err)
		assert.Equal(t, param.UpdateTaskResponse{}, response)
		assert.EqualError(t, err, expectedError.Error())
		mockRepo.AssertExpectations(t)
	})
}
