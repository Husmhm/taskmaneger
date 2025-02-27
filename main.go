package main

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"taskmaneger/adapter/redis"
	"taskmaneger/config"
	"taskmaneger/delivery/httpserver"
	"taskmaneger/repository/postgresql"
	"taskmaneger/repository/postgresql/postgersqltask"
	"taskmaneger/repository/postgresql/postgresqluser"
	"taskmaneger/service/authservice"
	"taskmaneger/service/taskservice"
	user "taskmaneger/service/userservice"
	"taskmaneger/validator"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Printf("cfg2: %+v\n", cfg)

	userSvc, authSvc, taskSvc, redisRepo := setupservice(cfg)

	server := httpserver.New(userSvc, taskSvc, authSvc, cfg.Auth, redisRepo)

	server.Start()
}

func setupservice(cfg config.Config) (user.Service, authservice.Service, taskservice.Service, redis.Adapter) {
	postgresqlRepo, err := postgresql.NewDB("host=localhost port=5435 user=admin dbname=task_manager_db sslmode=disable password=password123")
	if err != nil {
		fmt.Println(err)
	}
	if err := postgresqlRepo.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	authSvc := authservice.New(cfg.Auth)

	redisRepo := redis.New(cfg.Redis)

	userPostgresql := postgresqluser.New(postgresqlRepo)
	uV := validator.New(userPostgresql)
	userSvc := user.New(userPostgresql, uV, authSvc)

	taskPosrgresql := postgersqltask.New(postgresqlRepo)
	taskSvc := taskservice.New(taskPosrgresql)

	return userSvc, authSvc, taskSvc, redisRepo
}
