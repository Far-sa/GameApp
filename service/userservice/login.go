package userservice

import (
	"fmt"
	"game-app/param"
	"game-app/pkg/richerror"
)

func (s Service) Login(req param.LoginRequest) (param.LoginResponse, error) {
	const op = "userservice.Login"
	//! TODO --> it is better to use SOLID principle- imporove functionality for each task separately
	// check the phone number which is already exist
	// get user by phone number
	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if user.Password != getMD5Hash(req.Password) {
		return param.LoginResponse{}, fmt.Errorf("username/ password incorrect")
	}

	// create tokens
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	refreshToken, err := s.auth.RefreshAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	return param.LoginResponse{
		User:   param.UserInfo{ID: user.ID, Name: user.Name, PhoneNumber: user.PhoneNumber},
		Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
	}, nil
}
