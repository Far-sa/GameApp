package main

import (
	"fmt"
	"game-app/entity"
	"game-app/repository/mysql"
)

func main() {
	mysqlRepo := mysql.NewMYSQL()

	user, err := mysqlRepo.RegisterUser(entity.User{
		Name:        "tx",
		PhoneNumber: "0916",
	})
	if err != nil {
		fmt.Println(err)
	} else {

		fmt.Println(user)
	}

	isUnique, err := mysqlRepo.UniquenePhonenumber(user.PhoneNumber)
	if err != nil {
		fmt.Println("unique error", err)
	}

	fmt.Println("unique :", isUnique)
}
