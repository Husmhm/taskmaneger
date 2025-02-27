package config

import (
	"taskmaneger/adapter/redis"
	"taskmaneger/service/authservice"
)

type Config struct {
	Auth  authservice.Config `koanf:"auth"`
	Redis redis.Config       `koanf:"redis"`
}
