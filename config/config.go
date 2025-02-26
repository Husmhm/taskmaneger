package config

import "taskmaneger/service/authservice"

type Config struct {
	Auth authservice.Config `koanf:"auth"`
}
