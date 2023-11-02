package config

import (
	"game-app/adapter/redis"
	"game-app/repository/mysql"
	"game-app/service/authservice"
	"game-app/service/matchingservice"
	"time"
)

type Application struct {
	GracefullShutdownTimeout time.Duration `koanf:"gracefull_shutdown_timeout"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	Application Application            `koanf:"application`
	HTTPServer  HTTPServer             `koanf:"http_server"`
	Auth        authservice.Config     `koanf:"auth`
	Mysql       mysql.Config           `koanf:"mysql"`
	MatchingSvc matchingservice.Config `koanf:"matching_service`
	Redis       redis.Config           `koanf:"redis"`
}
