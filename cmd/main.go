package main

import (
	"fmt"
	"game-app/config"
	"game-app/delivery/httpserver"
	"game-app/repository/mysql"
	"game-app/repository/mysql/accessctl"
	"game-app/repository/mysql/mysqluser"
	"game-app/service/authorizationservice"
	"game-app/service/authservice"
	"game-app/service/backofficeuserservice"
	"game-app/service/matchingservice"
	"game-app/service/userservice"
	"game-app/validator/matchingvalidator"
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
	// TODO read config path from cmd
	cfg2 := config.Load("config.yml")
	fmt.Printf("cfg2 : %+v\n", cfg2)

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

	// TODO : add struct
	authSrv, userSrv, userValidator, backofficeUserSvc,
		authorizationSvc, matchingSvc, matchingV := setupServices(cfg)

	server := httpserver.New(
		cfg, authSrv, userSrv, userValidator, backofficeUserSvc,
		authorizationSvc, matchingSvc, matchingV)

	server.Serve()

}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service,
	uservalidator.Validator, backofficeuserservice.Service, authorizationservice.Service,
	matchingservice.Service, matchingvalidator.Validator,
) {
	authSrv := authservice.New(cfg.Auth)

	MysqlRepo := mysql.NewMYSQL(cfg.Mysql)

	userMysql := mysqluser.New(MysqlRepo)
	userSvc := userservice.New(authSrv, userMysql)

	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := accessctl.New(MysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	uV := uservalidator.New(userMysql)

	matchingV := matchingvalidator.New()
	matchingSvc := matchingservice.New(cfg.MatchingSvc)

	return authSrv, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingSvc, matchingV
}
