package redismatching

import (
	"context"
	"fmt"
	"game-app/entity"
	"game-app/pkg/richerror"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	WaitinglistPrefix = "waitinglist"
)

func (d DB) AddToWaitList(userID uint, category entity.Category) error {
	const op = "redismatching.AddToWaitList"

	_, err := d.adapter.Client().ZAdd(
		context.Background(),
		fmt.Sprintf("%s:%s", WaitinglistPrefix, category),
		redis.Z{
			Score:  float64(time.Now().UnixMicro()),
			Member: fmt.Sprintf("%d", userID),
		}).Result()

	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}
