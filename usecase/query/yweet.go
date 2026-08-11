package query

import (
	"context"
	"time"
)

type YweetQueryService interface {
	FindYweetByID(ctx context.Context, yweetID uint64) (*FindYweetByIDDTO, error)
}

type FindYweetByIDDTO struct {
	YweetID        uint64    `db:"yweet_id"`
	YweetContent   string    `db:"yweet_content"`
	YweetCreatedAt time.Time `db:"yweet_created_at"`

	UserID        uint64    `db:"user_id"`
	UserUsername  string    `db:"user_username"`
	UserCreatedAt time.Time `db:"user_created_at"`
}
