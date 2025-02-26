package postgresqluser

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	models "taskmaneger/model"
	"taskmaneger/pkg/errmsg"
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
func (d *DB) GetUserByPhoneNumber(phoneNumber string) (models.User, error) {

	var user models.User
	err := d.conn.Conn.Where("phone_number = ?", phoneNumber).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, fmt.Errorf(errmsg.ErrorMsgNotFound, err)
		}
		// TODO - log unexpected err for better observability
		return models.User{}, err
	}

	return user, nil
}
