package postgresqltask_test

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	models "taskmaneger/model"
	"taskmaneger/repository/postgresql/postgersqltask"

	"testing"
	"time"
)

func TestDBCreateTask(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	dbInstance := postgersqltask.New(db)

	t.Run("Create a normal task", func(t *testing.T) {
		newTask := models.Task{
			Title:       "Normal Task",
			Description: "This is a normal task",
			Status:      models.StatusInprogress,
			UserID:      1,
		}

		createdTask, err := dbInstance.CreateTask(newTask)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, createdTask.ID)
		assert.Equal(t, newTask.Title, createdTask.Title)
		assert.Equal(t, newTask.Description, createdTask.Description)
		assert.Equal(t, newTask.Status, createdTask.Status)
		assert.Equal(t, newTask.UserID, createdTask.UserID)
		assert.WithinDuration(t, time.Now(), createdTask.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), createdTask.UpdatedAt, time.Second)
	})

}

func TestGetTaskById(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	dbInstance := postgersqltask.New(db)

	newTask := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      models.StatusInprogress,
		UserID:      1,
	}

	createdTask, err := dbInstance.CreateTask(newTask)
	assert.NoError(t, err)

	t.Run("Get existing task by ID", func(t *testing.T) {
		task, err := dbInstance.GetTaskById(createdTask.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdTask.ID, task.ID)
		assert.Equal(t, createdTask.Title, task.Title)
		assert.Equal(t, createdTask.Description, task.Description)
		assert.Equal(t, createdTask.Status, task.Status)
		assert.Equal(t, createdTask.UserID, task.UserID)
	})

	t.Run("Get non-existing task by ID", func(t *testing.T) {
		nonExistingID := uint(999)
		_, err := dbInstance.GetTaskById(nonExistingID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "task not found")
	})
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	dbInstance := postgersqltask.New(db)

	newTask := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      models.StatusInprogress,
		UserID:      1,
	}

	createdTask, err := dbInstance.CreateTask(newTask)
	assert.NoError(t, err)

	t.Run("Update existing task", func(t *testing.T) {
		updatedTask := models.Task{
			ID:          createdTask.ID,
			Title:       "Updated Task",
			Description: "This task has been updated",
			Status:      models.StatusDone,
			UserID:      1,
		}

		result, err := dbInstance.UpdateTask(updatedTask)
		assert.NoError(t, err)
		assert.Equal(t, updatedTask.Title, result.Title)
		assert.Equal(t, updatedTask.Description, result.Description)
		assert.Equal(t, updatedTask.Status, result.Status)
	})

	t.Run("Update non-existing task", func(t *testing.T) {
		nonExistingTask := models.Task{
			ID:          99,
			Title:       "Non-existing Task",
			Description: "This task does not exist",
			Status:      models.StatusInprogress,
			UserID:      1,
		}

		result, err := dbInstance.UpdateTask(nonExistingTask)
		assert.Equal(t, result, models.Task{})
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	dbInstance := postgersqltask.New(db)

	newTask := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      models.StatusInprogress,
		UserID:      1,
	}

	createdTask, err := dbInstance.CreateTask(newTask)
	assert.NoError(t, err)

	t.Run("Delete existing task", func(t *testing.T) {
		err := dbInstance.DeleteTask(createdTask.ID)
		assert.NoError(t, err)

		var task models.Task
		err = db.Conn.First(&task, createdTask.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Delete non-existing task", func(t *testing.T) {
		nonExistingID := uint(999)
		err := dbInstance.DeleteTask(nonExistingID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestGetListOfTaskTitels(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	dbInstance := postgersqltask.New(db)

	userID := uint(1)

	tasks := []models.Task{
		{Title: "Task 1", Description: "Description 1", Status: models.StatusTodo, UserID: userID},
		{Title: "Task 2", Description: "Description 2", Status: models.StatusInprogress, UserID: userID},
		{Title: "Task 3", Description: "Description 3", Status: models.StatusDone, UserID: userID},
	}

	for _, task := range tasks {
		_, err := dbInstance.CreateTask(task)
		assert.NoError(t, err)
	}

	t.Run("Get list of task titles for existing user", func(t *testing.T) {
		titles, err := dbInstance.GetListOfTaskTitels(userID)
		assert.NoError(t, err)
		assert.Equal(t, []string{"Task 1", "Task 2", "Task 3"}, titles)
	})

	t.Run("Get list of task titles for non-existing user", func(t *testing.T) {
		nonExistingUserID := uint(999)
		titles, err := dbInstance.GetListOfTaskTitels(nonExistingUserID)
		assert.NoError(t, err)
		assert.Empty(t, titles)
	})
}
