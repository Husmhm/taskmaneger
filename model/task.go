package models

import (
	"time"
)

type TaskStatus uint

const (
	StatusTodo TaskStatus = iota
	StatusInprogress
	StatusDone
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      TaskStatus `gorm:"type:int;default:0" json:"status"`
	UserID      uint       `gorm:"not null" json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func TaskStatusIsValid(value uint) bool {
	return value <= uint(StatusDone)
}
