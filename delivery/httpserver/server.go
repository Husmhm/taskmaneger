package httpserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"taskmaneger/delivery/httpserver/userhandler"
	"taskmaneger/service/authservice"
	user "taskmaneger/service/userservice"
)

type Server struct {
	userhandler userhandler.Handler
	Router      *echo.Echo
}

func New(userSvc user.Service, authSvc authservice.Service) Server {
	return Server{
		userhandler: userhandler.New(userSvc, authSvc),
		Router:      echo.New(),
	}
}
func (s Server) Start() {
	s.userhandler.SetRoutes(s.Router)

	address := fmt.Sprintf("localhost:%d", 8088)
	fmt.Printf("star echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Printf("router start server error: %v\n", err)
	}
}
