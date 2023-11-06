package redismatching

import (
	"context"
	"fmt"
	"game-app/entity"
	"game-app/pkg/richerror"
	"game-app/pkg/timestamp"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

// TODO : add to config in usecase layer
const (
	WaitinglistPrefix = "waitinglist"
)

func (d DB) AddToWaitList(userID uint, category entity.Category) error {
	const op = "redismatching.AddToWaitList"

	_, err := d.adapter.Client().ZAdd(
		context.Background(),
		fmt.Sprintf("%s:%s", WaitinglistPrefix, category),
		redis.Z{
			Score:  float64(timestamp.Now()),
			Member: fmt.Sprintf("%d", userID),
		}).Result()

	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d DB) GetWaitListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error) {
	const op = "redismatching.GetWaitListByCategory"

	//d.adapter.Client().ZRangeWithScores()
	list, err := d.adapter.Client().ZRangeByScoreWithScores(ctx,
		getCategoryKey(category), &redis.ZRangeBy{
			//! convert int  to string
			Min:    fmt.Sprintf("%d", timestamp.Add(-2*time.Hour)),
			Max:    strconv.Itoa(int(timestamp.Now())),
			Offset: 0,
			Count:  0,
		}).Result()
	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	var result = make([]entity.WaitingMember, 0)

	for _, l := range list {
		userID, _ := strconv.Atoi(l.Member.(string))

		result = append(result, entity.WaitingMember{
			UserID:    uint(userID),
			Timestamp: int64(l.Score),
			Category:  category,
		})
	}
	return result, nil
}

func (d DB) RemoveUsersFromWaitingList(category entity.Category, userIDs []uint) {
	// TODO: add to config
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	members := make([]any, 0)
	for _, u := range userIDs {
		members = append(members, strconv.Itoa(int(u)))
	}

	numberOfRemovedMemberes, err := d.adapter.Client().ZRem(ctx, getCategoryKey(category), members...).Result()
	if err != nil {
		log.Errorf("remove users from waiting list : %v\n", err)
		// TODO: update metrics
	}

	log.Printf("%d items removed from %s", numberOfRemovedMemberes, getCategoryKey(category))

}

func getCategoryKey(category entity.Category) string {
	return fmt.Sprintf("%s:%s", WaitinglistPrefix, category)
}
