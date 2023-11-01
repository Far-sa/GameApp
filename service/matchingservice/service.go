package matchingservice

import (
	"game-app/entity"
	"game-app/param"
	"game-app/pkg/richerror"
	"time"
)

type Repo interface {
	AddToWaitList(userID uint, category entity.Category) error
}
type Config struct {
	WaitingTimeout time.Duration `koanf:"waiting_timeout"`
}

type Service struct {
	config Config
	repo   Repo
}

func New(config Config, repo Repo) Service {
	return Service{config: config, repo: repo}
}

func (s Service) AddToWaitingList(req param.AddToWaitingListRequest) (param.AddToWaitingListResponse, error) {
	const op = "matchingservice.AddToWaitingList"

	err := s.repo.AddToWaitList(req.UserID, req.Category)
	if err != nil {
		return param.AddToWaitingListResponse{}, richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected)
	}

	return param.AddToWaitingListResponse{}, nil
}
