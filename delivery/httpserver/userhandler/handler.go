package userhandler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"taskmaneger/param"
	user "taskmaneger/service/userservice"
)

type Handler struct {
	userSvc user.Service
}

func New(userSvc user.Service) Handler {
	return Handler{userSvc: userSvc}
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
