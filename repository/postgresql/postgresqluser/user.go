package postgresqluser

import (
	"context"
	"github.com/jinzhu/gorm"
	models "taskmaneger/model"
)

func (d *DB) Register(ctx context.Context, u models.User) (models.User, error) {

	if err := d.conn.Conn.Create(&u).Error; err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (d *DB) IsPhoneNumberUnique(ctx context.Context, phoneNumber string) (bool, error) {
	var user models.User
	err := d.conn.Conn.Where("phone_number = ?", phoneNumber).First(&user).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {

			return true, nil
		}

		return false, err
	}

	return false, nil
}
