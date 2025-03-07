package postgeresqluser_test

import (
	"github.com/stretchr/testify/assert"
	models "taskmaneger/model"
	"taskmaneger/repository/postgresql/postgresqluser"
	"testing"
)

func TestRegister(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	dbInstance := postgresqluser.New(db)

	t.Run("Register new user successfully", func(t *testing.T) {
		newUser := models.User{
			Name:        "testuser",
			PhoneNumber: "09131921277",
			Password:    "password123",
		}

		registeredUser, err := dbInstance.Register(newUser)
		assert.NoError(t, err)
		assert.NotEqual(t, uint(0), registeredUser.ID)
		assert.Equal(t, newUser.Name, registeredUser.Name)
		assert.Equal(t, newUser.PhoneNumber, registeredUser.PhoneNumber)
		assert.Equal(t, newUser.Password, registeredUser.Password)
	})

	t.Run("Register user with duplicate Phone number", func(t *testing.T) {
		newUser := models.User{
			Name:        "testuser2",
			PhoneNumber: "09131921277",
			Password:    "password123",
		}

		_, err := dbInstance.Register(newUser)
		assert.Error(t, err)
	})
}

func TestIsPhoneNumberUnique(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	dbInstance := postgresqluser.New(db)

	phoneNumber := "09131921288"

	t.Run("Phone number is unique", func(t *testing.T) {
		isUnique, err := dbInstance.IsPhoneNumberUnique(phoneNumber)
		assert.NoError(t, err)
		assert.True(t, isUnique)
	})

	t.Run("Phone number is not unique", func(t *testing.T) {
		newUser := models.User{
			Name:        "testuser",
			Password:    "password123",
			PhoneNumber: phoneNumber,
		}

		_, err := dbInstance.Register(newUser)
		assert.NoError(t, err)

		isUnique, err := dbInstance.IsPhoneNumberUnique(phoneNumber)
		assert.NoError(t, err)
		assert.False(t, isUnique)
	})

}

func TestGetUserByPhoneNumber(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	dbInstance := postgresqluser.New(db)

	phoneNumber := "09131921277"

	t.Run("Get existing user by phone number", func(t *testing.T) {
		newUser := models.User{
			Name:        "testuser",
			PhoneNumber: phoneNumber,
			Password:    "password123",
		}

		_, err := dbInstance.Register(newUser)
		assert.NoError(t, err)

		user, err := dbInstance.GetUserByPhoneNumber(phoneNumber)
		assert.NoError(t, err)
		assert.Equal(t, newUser.Name, user.Name)
		assert.Equal(t, newUser.PhoneNumber, user.PhoneNumber)
		assert.Equal(t, newUser.Password, user.Password)
	})

	t.Run("Get non-existing user by phone number", func(t *testing.T) {
		nonExistingPhoneNumber := "0000000000"
		_, err := dbInstance.GetUserByPhoneNumber(nonExistingPhoneNumber)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
