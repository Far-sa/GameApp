package main

import (
	"game-app/config"
	"game-app/delivery/httpserver"
	"game-app/repository/mysql"
	"game-app/service/authservice"
	"game-app/service/userservice"
	"game-app/validator/uservalidator"
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

	//* add migrator
	// mgr := migrator.New(cfg.Mysql)
	// mgr.Up()

	authSrv, userSrv, userValidator := setupServices(cfg)

	server := httpserver.New(cfg, authSrv, userSrv, userValidator)

	server.Serve()

}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator) {
	authSrv := authservice.New(cfg.Auth)

	MysqlRepo := mysql.NewMYSQL(cfg.Mysql)
	userSrv := userservice.New(authSrv, MysqlRepo)

	uV := uservalidator.New(MysqlRepo)

	return authSrv, userSrv, uV
}
