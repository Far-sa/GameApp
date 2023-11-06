package redis

import (
	"context"
	"game-app/entity"
	"time"
)

func (a Adapter) Publish(event entity.Event, payload string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := a.client.Publish(ctx, string(event), payload).Err(); err != nil {
		panic(err)
		// TODO: log error
	}

	//TODO: update metrics
}
