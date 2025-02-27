package middleware

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func RateLimiter(redisClient *redis.Client, limit int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := "rate_limit:" + ip

			currentCount, err := redisClient.Get(c.Request().Context(), key).Int()
			if err != nil && err != redis.Nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}

			if currentCount >= limit {
				return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "Too many requests"})
			}

			if err == redis.Nil {
				redisClient.Set(c.Request().Context(), key, 1, window)
			} else {
				redisClient.Incr(c.Request().Context(), key)
			}

			return next(c)
		}
	}
}
