package authservice

import (
	"game-app/entity"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ! interface pollution
// type AuthGenerator interface {
// 	CreateAccessToken(user entity.User) (string, error)
// 	RefreshAccessToken(user entity.User) (string, error)
// }

// type AuthParser interface {
// 	ParseToken(bearerToken string) (*Claims, error)
// }

type Config struct {
	SignKey               string
	AccessExpirationTime  time.Duration
	RefreshExpirationTime time.Duration
	AccessSubject         string
	RefreshSubject        string
}

type Service struct {
	config Config
}

func New(config Config) Service {
	return Service{config: config}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.AccessSubject, s.config.AccessExpirationTime)
}

func (s Service) RefreshAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.RefreshSubject, s.config.RefreshExpirationTime)
}

func (s Service) VerifyToken(bearerToken string) (*Claims, error) {

	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})

	if err != nil {
		return nil, err
	}

	// convert interface to conceret object
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil

	} else {
		return nil, err
	}
}

func (s Service) createToken(userID uint, subject string, expiresDuration time.Duration) (string, error) {
	// create a signer for rsa 256
	//t := jwt.New(jwt.GetSigningMethod("RS256"))
	// TODO replace with rsa 256 RS256

	// set our claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: subject,
			// set the expire time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresDuration)),
		},
		UserID: userID,
	}

	// TODO add sign method to config
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := accessToken.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
