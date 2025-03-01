package task_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"taskmaneger/param"
	"taskmaneger/service/taskservice"
	"testing"
	"time"
)

func TestService_List(t *testing.T) {
	userID := uint(1)
	validRequest := param.ListTaskTitlesRequest{}

	t.Run("Data found in cache", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		cacheKey := "task_titles:" + string(rune(userID))
		cachedResponse := param.ListTaskTitlesResponse{
			Tasks: []string{"Task 1", "Task 2"},
		}

		mockRedisRepo.On("Get", cacheKey, &param.ListTaskTitlesResponse{}).Run(func(args mock.Arguments) {
			dest := args.Get(1).(*param.ListTaskTitlesResponse)
			*dest = cachedResponse
		}).Return(nil)

		response, err := service.List(validRequest, userID)

		assert.NoError(t, err)
		assert.Equal(t, cachedResponse, response)
		mockRedisRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "GetListOfTaskTitels")
	})

	t.Run("Data not found in cache, fetch from database", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		cacheKey := "task_titles:" + string(rune(userID))
		titleList := []string{"Task 1", "Task 2"}
		dbResponse := param.ListTaskTitlesResponse{
			Tasks: titleList,
		}

		mockRedisRepo.On("Get", cacheKey, &param.ListTaskTitlesResponse{}).Return(errors.New("not found"))

		mockRepo.On("GetListOfTaskTitels", userID).Return(titleList, nil)

		mockRedisRepo.On("Set", cacheKey, dbResponse, 10*time.Minute).Return(nil)

		response, err := service.List(validRequest, userID)

		assert.NoError(t, err)
		assert.Equal(t, dbResponse, response)
		mockRedisRepo.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error fetching data from database", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		cacheKey := "task_titles:" + string(rune(userID))
		expectedError := errors.New("database error")

		mockRedisRepo.On("Get", cacheKey, &param.ListTaskTitlesResponse{}).Return(errors.New("not found"))

		mockRepo.On("GetListOfTaskTitels", userID).Return([]string{}, expectedError)

		response, err := service.List(validRequest, userID)

		assert.Error(t, err)
		assert.Equal(t, param.ListTaskTitlesResponse{}, response)
		assert.EqualError(t, err, expectedError.Error())
		mockRedisRepo.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error caching data", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockRedisRepo := new(MockRedisRepo)

		service := taskservice.New(mockRepo, mockRedisRepo)

		cacheKey := "task_titles:" + string(rune(userID))
		titleList := []string{"Task 1", "Task 2"}
		dbResponse := param.ListTaskTitlesResponse{
			Tasks: titleList,
		}
		cacheError := errors.New("cache error")

		mockRedisRepo.On("Get", cacheKey, &param.ListTaskTitlesResponse{}).Return(errors.New("not found"))

		mockRepo.On("GetListOfTaskTitels", userID).Return(titleList, nil)

		mockRedisRepo.On("Set", cacheKey, dbResponse, 10*time.Minute).Return(cacheError)

		response, err := service.List(validRequest, userID)

		assert.NoError(t, err)
		assert.Equal(t, dbResponse, response)
		mockRedisRepo.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
