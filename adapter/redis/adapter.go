package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
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

func New(config Config) Adapter {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(" %s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       0, // use default DB
	})

	return Adapter{client: rdb}
}

func (a Adapter) Client() *redis.Client {
	return a.client
}
