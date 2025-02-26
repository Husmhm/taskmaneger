package userhandler

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/users")
	// we can use skipper
	userGroup.POST("/register", h.userRegister)
	userGroup.POST("/login", h.userLogin)
}
