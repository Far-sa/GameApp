package matchingservice

import (
	"context"
	"fmt"
	"game-app/entity"
	"game-app/param"
	"game-app/pkg/richerror"
	"game-app/pkg/timestamp"
	"sync"
	"time"

	funk "github.com/thoas/go-funk"
)

// TODO : add context to all repo and use case methods
type Repo interface {
	AddToWaitList(userID uint, category entity.Category) error
	GetWaitListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
}

type PresenceClient interface {
	GetPresence(ctx context.Context, req param.GetPresenceRequest) (param.GetPresenceResponse, error)
}

type Config struct {
	WaitingTimeout time.Duration `koanf:"waiting_timeout"`
}

type Service struct {
	config         Config
	repo           Repo
	presenceClient PresenceClient
}

func New(config Config, repo Repo, presenceClient PresenceClient) Service {
	return Service{config: config, repo: repo, presenceClient: presenceClient}
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

func (s Service) MatchWaitedUsers(ctx context.Context, req param.MatchWaitedUsersRequest) (param.MatchWaitedUsersResponse, error) {
	const op = "matchingservice.MatchWaitedUsers"

	var wg sync.WaitGroup
	for _, category := range entity.CategoryList() {
		wg.Add(1)
		go s.match(ctx, category, &wg)
	}

	wg.Wait()
	return param.MatchWaitedUsersResponse{}, nil
}

func (s Service) match(ctx context.Context, category entity.Category, wg *sync.WaitGroup) {
	const op = "matchingservice.Match"

	defer wg.Done()

	list, err := s.repo.GetWaitListByCategory(ctx, category)
	if err != nil {
		//TODO : log err
		//TODO : update metrics
		return
	}

	userIDs := make([]uint, 0)
	for _, l := range list {
		userIDs = append(userIDs, l.UserID)
	}

	if len(userIDs) < 2 {
		return
	}

	//------->
	presenceList, err := s.presenceClient.GetPresence(ctx, param.GetPresenceRequest{UserID: userIDs})
	if err != nil {
		//TODO : log err
		//TODO : update metrics
		return
	}

	presenceUserIDs := make([]uint, len(list))
	for _, l := range presenceList.Items {
		presenceUserIDs = append(presenceUserIDs, l.UserID)
	}

	// TODO: merge presence list with user id list
	// consider the presence timestamp
	// remove user from waiting list if user's timestamp is older than <time>

	var finalList = make([]entity.WaitingMember, 0)
	for _, l := range list {
		if funk.ContainsUInt(presenceUserIDs, l.UserID) && l.Timestamp < timestamp.Add(-20*time.Second) {
			finalList = append(finalList, l)
		} else {
			// remove from list
		}
	}

	for i := 0; i < len(finalList)-1; i = i + 2 {

		mu := entity.MatchUsers{
			Category: category,
			UserIDs:  []uint{finalList[i].UserID, finalList[i+1].UserID},
		}
		// publish a new event for mu

		// remove mu users from waiting list
		fmt.Println("mu", mu)
	}
}
