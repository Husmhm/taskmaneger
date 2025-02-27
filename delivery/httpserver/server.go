package httpserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"taskmaneger/adapter/redis"
	"taskmaneger/delivery/httpserver/taskhandler"
	"taskmaneger/delivery/httpserver/userhandler"
	"taskmaneger/service/authservice"
	"taskmaneger/service/taskservice"
	user "taskmaneger/service/userservice"
)

type Server struct {
	userhandler userhandler.Handler
	taskhandler taskhandler.Handler
	Router      *echo.Echo
}

func New(userSvc user.Service, taskSvc taskservice.Service, authSvc authservice.Service, authCfg authservice.Config, redisRepo redis.Adapter) Server {
	return Server{
		userhandler: userhandler.New(userSvc),
		taskhandler: taskhandler.New(taskSvc, authSvc, authCfg, redisRepo),
		Router:      echo.New(),
	}
}
func (s Server) Start() {
	s.userhandler.SetRoutes(s.Router)
	s.taskhandler.SetRoutes(s.Router)

	address := fmt.Sprintf("localhost:%d", 8088)
	fmt.Printf("star echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Printf("router start server error: %v\n", err)
	}
}
