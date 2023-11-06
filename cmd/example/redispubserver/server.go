package main

import (
	"context"
	"fmt"
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/entity"
	"game-app/pkg/protobufencoder"
)

func main() {
	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)

	topic := "matching.users_matched"

	mu := entity.MatchUsers{
		Category: entity.SportCategory,
		UserIDs:  []uint{1, 4},
	}

	payload := protobufencoder.EncodeEvent(entity.MatchingUserEvent, mu)

	if err := redisAdapter.Client().Publish(context.Background(),
		topic, payload).Err(); err != nil {
		panic(fmt.Sprintf("publish error: %v", err))
	}
}
