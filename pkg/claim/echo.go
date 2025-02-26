package claim

import (
	"github.com/labstack/echo/v4"
	"taskmaneger/config"
	"taskmaneger/service/authservice"
)

func GetClaimsFromEchoContext(c echo.Context) *authservice.Claims {
	return c.Get(config.AuthMiddleWareContextKey).(*authservice.Claims)
}
