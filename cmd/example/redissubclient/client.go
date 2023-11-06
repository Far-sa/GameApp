package main

import (
	"context"
	"fmt"
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/pkg/protobufencoder"
)

func main() {

	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)

	topic := "matching.users_matched"

	subscriber := redisAdapter.Client().Subscribe(context.Background(), topic)

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}

		switch msg.Channel {
		case topic:
			processUsersMatchedEvent(msg.Channel, msg.Payload)
		default:
			fmt.Println("invaid topicg")
		}
	}
}

func processUsersMatchedEvent(topic string, data string) {

	mu := protobufencoder.DecodeMatchedUsersEvent(data)

	fmt.Println("received messages from" + topic + "topic.")
	fmt.Printf("matched users : %v\n", mu)
}
