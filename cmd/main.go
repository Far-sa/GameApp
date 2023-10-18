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

// func userProfileHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		fmt.Fprintf(w, `{"error":"invalid method"}`)
// 	}

// 	authSvc := authservice.New(JwtSignKey, AccessTokenSubject,
// 		RefreshTokenSubject, AccessTokenExpirationDuration, RefreshTokenExpirationDuration)

// 	authToken := r.Header.Get("Authorization")
// 	claims, err := authSvc.VerifyToken(authToken)
// 	if err != nil {
// 		fmt.Fprintf(w, "token is invalid")
// 	}

// 	mysqlRepo := mysql.NewMYSQL()
// 	userSvc := userservice.New(authSvc, mysqlRepo)

// 	resp, err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
// 	if err != nil {
// 		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

// 		return
// 	}

// 	data, err := json.Marshal(resp)
// 	if err != nil {
// 		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

// 		return
// 	}

// 	w.Write(data)
// }

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSrv := authservice.New(cfg.Auth)

	MysqlRepo := mysql.NewMYSQL(cfg.Mysql)
	userSrv := userservice.New(authSrv, MysqlRepo)

	return authSrv, userSrv
}
