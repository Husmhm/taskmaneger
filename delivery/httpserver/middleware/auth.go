package middleware

import (
	"fmt"
	mw "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	cfg "taskmaneger/config"
	"taskmaneger/service/authservice"
)

func Auth(service authservice.Service, config authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey:    cfg.AuthMiddleWareContextKey,
		SigningKey:    []byte(config.SignKey),
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.ParseToken(auth)
			if err != nil {
				return nil, err
			}
			fmt.Println(claims)
			return claims, nil
		},
	})
}
