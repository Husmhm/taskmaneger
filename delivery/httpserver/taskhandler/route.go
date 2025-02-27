package taskhandler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	mw "taskmaneger/delivery/httpserver/middleware"
	"time"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	tasksGroup := e.Group("/tasks")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(mw.Auth(h.authSvc, h.authCfg))
	e.Use(mw.RateLimiter(h.redisRepo.Client(), 10, time.Minute))

	tasksGroup.POST("/create", h.CreateTask)
	tasksGroup.GET("/:id/view", h.ReadTask)
	tasksGroup.PUT("/:id/update", h.UpdateTask)
	tasksGroup.DELETE("/:id/delete", h.DeleteTask)
}
