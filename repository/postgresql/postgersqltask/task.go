package postgersqltask

import (
	"fmt"
	"github.com/jinzhu/gorm"
	models "taskmaneger/model"
)

func (d DB) CreateTask(task models.Task) (models.Task, error) {
	if err := d.conn.Conn.Create(&task).Error; err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (d DB) GetTaskById(id uint) (models.Task, error) {
	var task models.Task
	err := d.conn.Conn.First(&task, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Task{}, fmt.Errorf("task not found")
		}
		return models.Task{}, err
	}
	return task, nil
}

func (d DB) UpdateTask(task models.Task) (models.Task, error) {

	result := d.conn.Conn.Model(task).Updates(task)
	if result.Error != nil {
		return models.Task{}, result.Error
	}
	return task, nil
}

func (d DB) DeleteTask(id uint) error {
	result := d.conn.Conn.Where("id = ?", id).Delete(&models.Task{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
