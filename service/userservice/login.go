package userservice

import (
	"fmt"
	"game-app/dto"
	"game-app/pkg/richerror"
)

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
