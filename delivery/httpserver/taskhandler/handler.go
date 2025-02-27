package taskhandler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"taskmaneger/adapter/redis"
	"taskmaneger/param"
	"taskmaneger/pkg/claim"
	"taskmaneger/service/authservice"
	"taskmaneger/service/taskservice"
)

type Handler struct {
	taskSvc   taskservice.Service
	authSvc   authservice.Service
	authCfg   authservice.Config
	redisRepo redis.Adapter
}

func New(taskSvc taskservice.Service, authSvc authservice.Service, authCfg authservice.Config, redisRepo redis.Adapter) Handler {
	return Handler{taskSvc: taskSvc, authSvc: authSvc, authCfg: authCfg, redisRepo: redisRepo}
}

func (h Handler) CreateTask(c echo.Context) error {
	var req param.CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims := claim.GetClaimsFromEchoContext(c)

	resp, err := h.taskSvc.Create(req, claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h Handler) ReadTask(c echo.Context) error {
	var req param.ReadTaskRquest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims := claim.GetClaimsFromEchoContext(c)

	resp, err := h.taskSvc.Read(req, claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h Handler) UpdateTask(c echo.Context) error {
	var req param.UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	fmt.Println(req.TaskId)
	claims := claim.GetClaimsFromEchoContext(c)

	resp, err := h.taskSvc.Update(req, claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h Handler) DeleteTask(c echo.Context) error {
	var req param.DeleteTaskRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims := claim.GetClaimsFromEchoContext(c)

	resp, err := h.taskSvc.Delete(req, claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
