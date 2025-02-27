package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

type Adapter struct {
	client *redis.Client
}

func New(cfg Config) Adapter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})
	return Adapter{client: rdb}
}

func (a Adapter) Client() *redis.Client {
	return a.client
}
