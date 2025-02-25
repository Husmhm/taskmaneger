package models

import "gorm.io/gorm"

type taskStatus uint

const (
	StatusTodo taskStatus = iota
	StatusInprogress
	StatusDone
)

type Task struct {
	gorm.Model
	Title       string
	Description string
	Status      taskStatus // To-Do, In Progress, Done
	UserID      uint
}
