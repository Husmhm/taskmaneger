package taskhandler

import (
	"github.com/labstack/echo/v4"
	"taskmaneger/delivery/httpserver/middleware"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	tasksGroup := e.Group("/tasks")

	tasksGroup.POST("/create", h.CreateTask, middleware.Auth(h.authSvc, h.authCfg))
	tasksGroup.GET("/:id/view", h.ReadTask, middleware.Auth(h.authSvc, h.authCfg))
	tasksGroup.PUT("/:id/update", h.UpdateTask, middleware.Auth(h.authSvc, h.authCfg))
	tasksGroup.DELETE("/:id/delete", h.DeleteTask, middleware.Auth(h.authSvc, h.authCfg))
}
