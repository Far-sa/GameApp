package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"game-app/dto"
	"game-app/entity"
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

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {

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
	}

	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	//* return Created User
	return dto.RegisterResponse{User: dto.UserInfo{
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

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "userservice.Login"
	//! TODO --> it is better to use SOLID principle- imporove functionality for each task separately
	// check the phone number which is already exist
	// get user by phone number
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return dto.LoginResponse{}, fmt.Errorf("record not found %w", err)
	}

	if user.Password != getMD5Hash(req.Password) {
		return dto.LoginResponse{}, fmt.Errorf("username/ password incorrect")
	}

	// create tokens
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	refreshToken, err := s.auth.RefreshAccessToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return dto.LoginResponse{
		User:   dto.UserInfo{ID: user.ID, Name: user.Name, PhoneNumber: user.PhoneNumber},
		Tokens: dto.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
	}, nil
}

func (s Service) Profile(req dto.ProfileRequest) (dto.ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserById(uint64(req.UserID))
	if err != nil {
		//!  - use rich error to develope error handling from different layers
		return dto.ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})

		//return ProfileResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	return dto.ProfileResponse{Name: user.Name}, nil
}

func getMD5Hash(test string) string {
	hash := md5.Sum([]byte(test))
	return hex.EncodeToString(hash[:])
}
