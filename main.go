package main

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"taskmaneger/config"
	"taskmaneger/delivery/httpserver"
	"taskmaneger/repository/postgresql"
	"taskmaneger/repository/postgresql/postgresqluser"
	"taskmaneger/service/authservice"
	user "taskmaneger/service/userservice"
	"taskmaneger/validator"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Printf("cfg2: %+v\n", cfg)

	userSvc, authSvc := setupservice(cfg)
	server := httpserver.New(userSvc, authSvc)

	server.Start()
}

func setupservice(cfg config.Config) (user.Service, authservice.Service) {
	postgresqlRepo, err := postgresql.NewDB("host=localhost port=5435 user=admin dbname=task_manager_db sslmode=disable password=password123")
	if err != nil {
		fmt.Println(err)
	}
	if err := postgresqlRepo.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	authSvc := authservice.New(cfg.Auth)

	userPostgresql := postgresqluser.New(postgresqlRepo)
	uV := validator.New(userPostgresql)
	userSvc := user.New(userPostgresql, uV, authSvc)

	return userSvc, authSvc
}
