package repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"taskmaneger/model"
)

type DB struct {
	Conn *gorm.DB
}

func NewDB(connectionString string) (*DB, error) {
	db, err := gorm.Open("postgresql", connectionString)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	return &DB{Conn: db}, nil
}

func (d *DB) Close() error {
	return d.Conn.Close()
}

func (d *DB) AutoMigrate() error {
	err := d.Conn.AutoMigrate(
		&models.User{},
		&models.Task{},
	).Error

	if err != nil {
		return err
	}

	log.Println("Migration completed successfully!")
	return nil
}
