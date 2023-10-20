package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"game-app/entity"
	"game-app/pkg/phonenumber"
	"game-app/pkg/richerror"
)

type Repository interface {
	UniquenePhonenumber(phoneNumer string) (bool, error)
	RegisterUser(user entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserById(userID uint64) (entity.User, error)
}

// ! interface composibility -->  helpul for test
type AuthGeneratorService interface {
	CreateAccessToken(user entity.User) (string, error)
	RefreshAccessToken(user entity.User) (string, error)
}

type Service struct {
	// use auth service as interface
	auth AuthGeneratorService
	repo Repository
}

func New(authGenerator AuthGeneratorService, repo Repository) Service {
	return Service{auth: authGenerator, repo: repo}
}

// * DTOs
type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserInfo struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

type RegisterResponse struct {
	User UserInfo `json:"user"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {

	//? TODO --> verify phone number with verification code

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

	//? TODO use regex to validate password
	//* validate Password
	if len(req.Password) < 4 {
		return RegisterResponse{}, fmt.Errorf("password should be at least 4 characters")
	}

	//* third party bcrypt
	// bcrypt.GenerateFromPassword(pass, 0)

	//* save user in DB
	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    getMD5Hash(req.Password),
	}

	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	//* return Created User
	return RegisterResponse{User: UserInfo{
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

// * DTOs
type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type Tokens struct {
	// token
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User   UserInfo `json:"user"`
	Tokens Tokens   `json:"token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	const op = "userservice.Login"
	//! TODO --> it is better to use SOLID principle- imporove functionality for each task separately
	// check the phone number which is already exist
	// get user by phone number
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, richerror.New(op).WithErr(err).WithMessage("invalid info")
	}
	fmt.Printf("user found by phone: %+v", user)

	if !exist {
		return LoginResponse{}, fmt.Errorf("record not found %w", err)
	}

	if user.Password != getMD5Hash(req.Password) {
		return LoginResponse{}, fmt.Errorf("username/ password incorrect")
	}

	// create tokens
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	refreshToken, err := s.auth.RefreshAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return LoginResponse{
		User:   UserInfo{ID: user.ID, Name: user.Name, PhoneNumber: user.PhoneNumber},
		Tokens: Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
	}, nil
}

// * DTOs
type ProfileRequest struct {
	UserID uint `json:"user_id"`
}
type ProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserById(uint64(req.UserID))
	if err != nil {
		//!  - use rich error to develope error handling from different layers
		return ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})

		//return ProfileResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	return ProfileResponse{Name: user.Name}, nil
}

func getMD5Hash(test string) string {
	hash := md5.Sum([]byte(test))
	return hex.EncodeToString(hash[:])
}
