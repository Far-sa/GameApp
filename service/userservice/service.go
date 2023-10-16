package userservice

import (
	"fmt"
	"game-app/entity"
	"game-app/pkg/phonenumber"
)

type Repository interface {
	UniquenePhonenumber(phoneNumer string) (bool, error)
	RegisterUser(user entity.User) (entity.User, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

// * DTOs
type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type RegisterResponse struct {
	User entity.User `json:"user"`
}

type LoginRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type LoginResponse struct {
	// token
}

//* ---------------

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {

	//TODO --> verify phone number with verification code

	//* validate phone number && uniqueness
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	if isUnigue, err := s.repo.UniquenePhonenumber(req.PhoneNumber); err != nil || !isUnigue {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
		}
		if !isUnigue {
			return RegisterResponse{}, fmt.Errorf("phone number is not unique")
		}
	}

	//* validate  name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name should be at least 3 characters")
	}

	//* save user in DB
	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
	}

	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	//* return Created User
	return RegisterResponse{User: createdUser}, nil
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	panic("un implemented")
}
