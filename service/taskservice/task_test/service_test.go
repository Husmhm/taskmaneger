package task_test

import (
	"github.com/stretchr/testify/mock"
	models "taskmaneger/model"
	"time"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateTask(task models.Task) (models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockRepository) GetTaskById(taskId uint) (models.Task, error) {
	args := m.Called(taskId)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockRepository) UpdateTask(task models.Task) (models.Task, error) {
	args := m.Called(task)
	return args.Get(0).(models.Task), args.Error(1)
}

func (m *MockRepository) DeleteTask(taskId uint) error {
	args := m.Called(taskId)
	return args.Error(0)
}

func (m *MockRepository) GetListOfTaskTitels(userID uint) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

type MockRedisRepo struct {
	mock.Mock
}

func (m *MockRedisRepo) Get(key string, dest interface{}) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *MockRedisRepo) Set(key string, value interface{}, expireTime time.Duration) error {
	args := m.Called(key, value, expireTime)
	return args.Error(0)
}
