package matchingservice

import (
	"context"
	"game-app/entity"
	"game-app/param"
	"game-app/pkg/protobufencoder"
	"game-app/pkg/richerror"
	"game-app/pkg/timestamp"
	"sync"
	"time"
)

type Publisher interface {
	Publish(event entity.Event, payload string)
}

// TODO : add context to all repo and use case methods if necessary
type Repo interface {
	AddToWaitList(userID uint, category entity.Category) error
	GetWaitListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
	RemoveUsersFromWaitingList(category entity.Category, userIDs []uint)
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
	pub            Publisher
}

func New(config Config, repo Repo, presenceClient PresenceClient, pub Publisher) Service {
	return Service{config: config, repo: repo, presenceClient: presenceClient, pub: pub}
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

	//-------> rpc call
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

	var toBeRemovedUsers = make([]uint, 0)

	var finalList = make([]entity.WaitingMember, 0)

	for _, l := range list {
		lastOnlineTimestamp, ok := getPresenceItems(presenceList, l.UserID)
		if ok && lastOnlineTimestamp > timestamp.Add(-20*time.Second) &&
			l.Timestamp > timestamp.Add(-120*time.Second) {
			finalList = append(finalList, l)
		} else {
			toBeRemovedUsers = append(toBeRemovedUsers, l.UserID)
		}
	}

	//* no need to have main ctx and error -
	go s.repo.RemoveUsersFromWaitingList(category, toBeRemovedUsers)

	matchedUsersToBeRemoved := make([]uint, 0)
	for i := 0; i < len(finalList)-1; i = i + 2 {

		mu := entity.MatchUsers{
			Category: category,
			UserIDs:  []uint{finalList[i].UserID, finalList[i+1].UserID},
		}
		// publish a new event for mu
		go s.pub.Publish(entity.MatchingUserEvent,
			protobufencoder.EncodeEvent(entity.MatchingUserEvent, mu))

		// remove mu users from waiting list
		matchedUsersToBeRemoved = append(matchedUsersToBeRemoved, mu.UserIDs...)
	}

	go s.repo.RemoveUsersFromWaitingList(category, matchedUsersToBeRemoved)
}

// * custom fn which is act as map struct in entity
func getPresenceItems(presenceList param.GetPresenceResponse, userID uint) (int64, bool) {
	for _, item := range presenceList.Items {
		if item.UserID == userID {
			return item.Timestamp, true
		}
	}

	return 0, false
}
