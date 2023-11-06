package broker

import "game-app/entity"

type Publisher interface {
	Publish(event entity.Event, payload string)
}
