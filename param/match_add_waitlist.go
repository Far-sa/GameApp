package param

import (
	"game-app/entity"
	"time"
)

type AddToWaitingListRequest struct {
	UserID   uint `json:"user_id"`
	Category entity.Category
}

type AddToWaitingListResponse struct {
	Timeout time.Duration `json:"timeout,omitempty"`
}
