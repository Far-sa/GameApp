package presenceservice

import (
	"context"
	"fmt"
	"game-app/param"
	"game-app/pkg/richerror"
	"time"
)

type Config struct {
	ExpirationTiem time.Duration `koanf:"expiration_time"`
	Prefix         string        `koanf:"prefix"`
}

type Repository interface {
	Upsert(ctx context.Context, key string, timestamp int64, expTime time.Duration) error
	GetPresence(ctx context.Context, key string, userIDs []uint) (map[uint]int64, error)
}

type Service struct {
	config Config
	repo   Repository
}

func New(config Config, repo Repository) Service {
	return Service{config: config, repo: repo}
}

func (s Service) UpsertPresence(ctx context.Context, req param.UpsertPresenceRequest) (param.UpsertPresenceResponse, error) {
	const op = "presenceservice.UpsertPresence"

	err := s.repo.Upsert(ctx, fmt.Sprintf("%s:%d",
		s.config.Prefix, req.UserID),
		req.Timestamp, s.config.ExpirationTiem)
	if err != nil {
		fmt.Println("UpsertPresence err2 :", err.Error())
		return param.UpsertPresenceResponse{}, richerror.New(op).WithErr(err)
	}

	return param.UpsertPresenceResponse{}, nil
}

func (s Service) GetPresence(ctx context.Context, req param.GetPresenceRequest) (param.GetPresenceResponse, error) {
	fmt.Println("req:", req)

	// TODO: not implemented
	return param.GetPresenceResponse{Items: []param.GetPresenceItem{
		{UserID: 1, Timestamp: 6541},
		{UserID: 2, Timestamp: 687513},
	}}, nil
}
