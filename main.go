package main

import (
	"context"
	"fmt"
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/delivery/httpserver"
	"game-app/repository/mysql"
	"game-app/repository/mysql/accessctl"
	"game-app/repository/mysql/mysqluser"
	"game-app/repository/redis/redismatching"
	"game-app/schedular"
	"game-app/service/authorizationservice"
	"game-app/service/authservice"
	"game-app/service/backofficeuserservice"
	"game-app/service/matchingservice"
	"game-app/service/userservice"
	"game-app/validator/matchingvalidator"
	"game-app/validator/uservalidator"
	"os"
	"os/signal"
	"time"
)

func main() {
	// TODO read config path from cmd
	cfg := config.Load("config.yml")
	fmt.Printf("cfg : %+v\n", cfg)

	//* add migrator
	// mgr := migrator.New(cfg.Mysql)
	// mgr.Up()

	// TODO : add struct
	authSrv, userSrv, userValidator, backofficeUserSvc,
		authorizationSvc, matchingSvc, matchingV := setupServices(cfg)

	server := httpserver.New(
		cfg, authSrv, userSrv, userValidator, backofficeUserSvc,
		authorizationSvc, matchingSvc, matchingV)

	go func() {
		server.Serve()
	}()

	done := make(chan bool)
	go func() {
		sch := schedular.New()
		sch.Start(done)
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit
	fmt.Println("received interrupt signal,shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Application.GracefullShutdownTimeout)
	defer cancel()

	if err := server.Router.Shutdown(ctx); err != nil {
		fmt.Println("httpServer shutdown error:", err)
	}

	done <- true
	time.Sleep(cfg.Application.GracefullShutdownTimeout)

	// TODO : context doesn't wait for schedular to finish its jobs
	<-ctx.Done()

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

	redisAdapter := redis.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdapter)
	matchingSvc := matchingservice.New(cfg.MatchingSvc, matchingRepo)

	return authSrv, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingSvc, matchingV
}
