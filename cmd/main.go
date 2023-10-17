package main

import (
	"encoding/json"
	"fmt"
	"game-app/repository/mysql"
	"game-app/service/authservice"
	"game-app/service/userservice"
	"io"
	"log"
	"net/http"
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

	mux := http.NewServeMux()

	mux.HandleFunc("/health-check", healthCheckHandler)
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/users/login", userLoginHandler)
	mux.HandleFunc("/users/profile", userProfileHandler)

	log.Println("Server started on Port : 8000... ")
	server := http.Server{Addr: ":8000", Handler: mux}
	log.Fatal(server.ListenAndServe())
}

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, `{"error":"invalid method"}`)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	var req userservice.RegisterRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	authSvc := authservice.New(JwtSignKey, AccessTokenSubject,
		RefreshTokenSubject, AccessTokenExpirationDuration, RefreshTokenExpirationDuration)

	mysqlRepo := mysql.NewMYSQL()
	userSvc := userservice.New(authSvc, mysqlRepo)

	_, err = userSvc.Register(req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	w.Write([]byte(`{"message":"user created successfully"}`))

}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, `{"error":"invalid method"}`)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	var req userservice.LoginRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	authSvc := authservice.New(JwtSignKey, AccessTokenSubject,
		RefreshTokenSubject, AccessTokenExpirationDuration, RefreshTokenExpirationDuration)

	mysqlRepo := mysql.NewMYSQL()
	userSvc := userservice.New(authSvc, mysqlRepo)

	resp, err := userSvc.Login(req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	data, err = json.Marshal(resp)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	w.Write(data)
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprintf(w, `{"error":"invalid method"}`)
	}

	authSvc := authservice.New(JwtSignKey, AccessTokenSubject,
		RefreshTokenSubject, AccessTokenExpirationDuration, RefreshTokenExpirationDuration)

	authToken := r.Header.Get("Authorization")
	claims, err := authSvc.ParseToken(authToken)
	if err != nil {
		fmt.Fprintf(w, "token is invalid")
	}

	mysqlRepo := mysql.NewMYSQL()
	userSvc := userservice.New(authSvc, mysqlRepo)

	resp, err := userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	w.Write(data)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"everything is ok"}`)
}
