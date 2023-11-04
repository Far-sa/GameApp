package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/contract/golang/matching"
	"game-app/entity"
	"game-app/pkg/slice"

	"google.golang.org/protobuf/proto"
)

func main() {
	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)

	topic := "matching.users_matched"
	mu := entity.MatchUsers{
		Category: entity.SportCategory,
		UserIDs:  []uint{1, 4},
	}

	pbMu := matching.MatchUsers{
		Category: string(mu.Category),
		UserIds:  slice.MapFromUintToUint64(mu.UserIDs),
	}

	//! serialize payload to pass over network
	payload, err := proto.Marshal(&pbMu)
	if err != nil {
		panic(err)
	}

	//! saftey of payload
	payloadStr := base64.StdEncoding.EncodeToString(payload)

	if err := redisAdapter.Client().Publish(context.Background(),
		topic, payloadStr).Err(); err != nil {
		panic(fmt.Sprintf("publish error: %v", err))
	}
}
