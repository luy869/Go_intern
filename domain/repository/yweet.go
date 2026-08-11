package repository

import (
	"context"
	"yatter-backend-go/app/domain/object/yweet"
)

type Yweet interface {
	Insert(ctx context.Context, pendingYweet *yweet.PendingYweet) (*yweet.Yweet, error)
}
