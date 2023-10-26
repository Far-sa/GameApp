package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"game-app/entity"
)

type Repository interface {
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

func getMD5Hash(test string) string {
	hash := md5.Sum([]byte(test))
	return hex.EncodeToString(hash[:])
}
