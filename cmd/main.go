package main

import (
	"encoding/json"
	"fmt"
	"game-app/repository/mysql"
	"game-app/service/userservice"
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/health-check", healthCheckHandler)
	http.HandleFunc("/users/register", userRegisterHandler)

	log.Println("Server started on Port : 8000... ")
	http.ListenAndServe(":8000", nil)

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

	mysqlRepo := mysql.NewMYSQL()
	userSvc := userservice.New(mysqlRepo)

	_, err = userSvc.Register(req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"`, err.Error())))

		return
	}

	w.Write([]byte(`{"message":"user created successfully"}`))

}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"everything is ok"}`)
}