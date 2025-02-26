package taskservice

import (
	"fmt"
	models "taskmaneger/model"
	"taskmaneger/param"
	"time"
)

type Repository interface {
	CreateTask(task models.Task) (models.Task, error)
	GetTaskById(taskId uint) (models.Task, error)
	UpdateTask(task models.Task) (models.Task, error)
	DeleteTask(taskId uint) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Create(req param.CreateTaskRequest, userID uint) (param.CreateTaskResponse, error) {
	Task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      0,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	fmt.Println(Task)

	createdTask, err := s.repo.CreateTask(Task)
	if err != nil {
		return param.CreateTaskResponse{}, err
	}
	fmt.Println(createdTask)
	response := param.CreateTaskResponse{
		Id:     createdTask.ID,
		Title:  createdTask.Title,
		Status: createdTask.Status,
		UserID: createdTask.UserID,
	}
	return response, nil
}

func (s Service) Read(req param.ReadTaskRquest, userID uint) (param.ReadTaskResponse, error) {
	task, err := s.repo.GetTaskById(req.TaskId)
	if err != nil {
		return param.ReadTaskResponse{}, err
	}
	if task.UserID != userID {
		return param.ReadTaskResponse{}, fmt.Errorf("not authorized")
	}
	resp := param.ReadTaskResponse{Task: task}
	return resp, nil
}

func (s Service) Update(req param.UpdateTaskRequest, userID uint) (param.UpdateTaskResponse, error) {
	task, err := s.repo.GetTaskById(req.TaskId)
	if err != nil {
		return param.UpdateTaskResponse{}, err
	}
	if task.UserID != userID {
		return param.UpdateTaskResponse{}, fmt.Errorf("not authorized")
	}
	if !models.TaskStatusIsValid(uint(req.Status)) {
		return param.UpdateTaskResponse{}, fmt.Errorf("invalid status")
	}

	task = models.Task{
		ID:          req.TaskId,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		UserID:      userID,
		UpdatedAt:   time.Now(),
	}
	fmt.Println(task.ID)
	updatedTask, err := s.repo.UpdateTask(task)
	if err != nil {
		return param.UpdateTaskResponse{}, err
	}
	fmt.Println(updatedTask.ID)
	return param.UpdateTaskResponse{updatedTask}, nil
}

func (s Service) Delete(req param.DeleteTaskRequest, userID uint) (param.DeleteTaskResponse, error) {
	task, err := s.repo.GetTaskById(req.TaskId)
	if err != nil {
		return param.DeleteTaskResponse{}, err
	}
	if task.UserID != userID {
		return param.DeleteTaskResponse{}, fmt.Errorf("not authorized")
	}

	deleteErr := s.repo.DeleteTask(req.TaskId)

	if deleteErr != nil {
		return param.DeleteTaskResponse{}, deleteErr
	}

	return param.DeleteTaskResponse{"successfully delete"}, nil
}
