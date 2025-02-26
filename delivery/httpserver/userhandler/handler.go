package userhandler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"taskmaneger/param"
	"taskmaneger/service/authservice"
	user "taskmaneger/service/userservice"
)

type Handler struct {
	authSvc authservice.Service
	userSvc user.Service
}

func New(userSvc user.Service, authSvc authservice.Service) Handler {
	return Handler{userSvc: userSvc, authSvc: authSvc}
}

func (h Handler) userRegister(c echo.Context) error {
	var req param.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, err := h.userSvc.Register(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h Handler) userLogin(c echo.Context) error {
	var req param.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, err := h.userSvc.Login(req)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)

}
