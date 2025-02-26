package param

import models "taskmaneger/model"

type CreateTaskRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
	UserID      uint              `json:"user_id"`
}

type CreateTaskResponse struct {
	Id     uint              `json:"id"`
	Title  string            `json:"title"`
	Status models.TaskStatus `json:"status"`
	UserID uint              `json:"user_id"`
}

type ReadTaskRquest struct {
	TaskId uint `json:"-" param:"id"`
}

type ReadTaskResponse struct {
	Task models.Task `json:"task"`
}

type UpdateTaskRequest struct {
	TaskId      uint              `json:"-" param:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
}
type UpdateTaskResponse struct {
	Task models.Task `json:"task"`
}

type DeleteTaskRequest struct {
	TaskId uint `json:"-" param:"id"`
}
type DeleteTaskResponse struct {
	Message string `json:"message"`
}
