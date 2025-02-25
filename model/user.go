package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PhoneNumber string `gorm:"unique" json:"phone_number"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}
