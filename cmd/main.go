package main

import (
	"game-app/config"
	"game-app/delivery/httpserver"
	"game-app/repository/mysql"
	"game-app/service/authservice"
	"game-app/service/userservice"
	"time"
)

const (
	JwtSignKey                     = "jwt-secret"
	AccessTokenSubject             = "at"
	RefreshTokenSubject            = "rt"
	AccessTokenExpirationDuration  = time.Hour * 24
	RefreshTokenExpirationDuration = time.Hour * 24 * 7
)

func main() {

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8000},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
			AccessExpirationTime:  AccessTokenExpirationDuration,
			RefreshExpirationTime: RefreshTokenExpirationDuration,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
		},

		Mysql: mysql.Config{
			Username: "root",
			Password: "password",
			Host:     "localhost",
			Port:     3306,
			DbName:   "gamedb",
		},
	}

	authSrv, userSrv := setupServices(cfg)

	server := httpserver.New(cfg, authSrv, userSrv)

	server.Serve()

}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSrv := authservice.New(cfg.Auth)

	MysqlRepo := mysql.NewMYSQL(cfg.Mysql)
	userSrv := userservice.New(authSrv, MysqlRepo)

	return authSrv, userSrv
}
