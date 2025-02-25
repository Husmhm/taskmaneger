package main

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"taskmaneger/delivery/httpserver"
	"taskmaneger/repository/postgresql"
	"taskmaneger/repository/postgresql/postgresqluser"
	user "taskmaneger/service/userservice"
	"taskmaneger/validator"
)

func main() {
	userSvc := setupservice()
	server := httpserver.New(userSvc)

	server.Start()
}

func setupservice() user.Service {
	postgresqlRepo, err := postgresql.NewDB("host=localhost port=5434 user=admin dbname=task_manager_db sslmode=disable password=password123")
	if err != nil {
		fmt.Println(err)
	}
	if err := postgresqlRepo.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	userPostgresql := postgresqluser.New(postgresqlRepo)
	uV := validator.New(userPostgresql)
	userSvc := user.New(userPostgresql, uV)

	return userSvc
}
