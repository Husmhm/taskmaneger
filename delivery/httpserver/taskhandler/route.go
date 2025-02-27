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

	e.Use(mw.RateLimiter(h.redisRepo.Client(), 10, time.Minute))

	tasksGroup.POST("/create", h.CreateTask, mw.Auth(h.authSvc, h.authCfg))
	tasksGroup.GET("/:id/view", h.ReadTask, mw.Auth(h.authSvc, h.authCfg))
	tasksGroup.PUT("/:id/update", h.UpdateTask, mw.Auth(h.authSvc, h.authCfg))
	tasksGroup.DELETE("/:id/delete", h.DeleteTask, mw.Auth(h.authSvc, h.authCfg))
	tasksGroup.GET("/", h.ListTitleTasks, mw.Auth(h.authSvc, h.authCfg))
}
