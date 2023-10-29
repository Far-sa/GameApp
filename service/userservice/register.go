package userservice

import (
	"fmt"
	"game-app/entity"
	"game-app/param"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {

	//? TODO --> verify phone number with verification code

	//* --> assign reqiest validation in other service

	//* third party bcrypt
	// bcrypt.GenerateFromPassword(pass, 0)

	//* save user in DB
	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    getMD5Hash(req.Password),
		Role:        entity.UserRole,
	}

	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	//* return Created User
	return param.RegisterResponse{User: param.UserInfo{
		ID:          createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}}, nil

	//! anonymous struct
	// resp2 := RegisterResponse{User: struct {
	// 	ID          uint   `json:"id"`
	// 	Name        string `json:"name"`
	// 	PhoneNumber string `json:"phone_number`
	// }{ID: createdUser.ID, Name: createdUser.Name, PhoneNumber: createdUser.PhoneNumber}}

	//return resp2,nil
}
