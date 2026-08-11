package infra

import (
	"context"
	"time"
	"yatter-backend-go/app/domain/object/yweet"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/infra/transaction"
)

var _ repository.Yweet = (*YweetRepoImpl)(nil)

type YweetRepoImpl struct{}

type insertedYweetDTO struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

func NewYweetRepository() *YweetRepoImpl {
	return &YweetRepoImpl{}
}

func (ur *YweetRepoImpl) Insert(ctx context.Context, pendingYweet *yweet.PendingYweet) (*yweet.Yweet, error) {
	tx, err := transaction.FetchTransaction(ctx)
	if err != nil {
		return nil, err
	}

	insertResult, err := tx.ExecContext(
		ctx,
		`INSERT INTO yweet (user_id, content) VALUES (?, ?)`,
		pendingYweet.UserID(),
		pendingYweet.Content(),
	)
	if err != nil {
		return nil, err
	}

	yweetID, err := insertResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	var insertedYweetDTO insertedYweetDTO
	err = tx.GetContext(ctx, &insertedYweetDTO, `SELECT id, user_id, content, created_at FROM yweet WHERE id = ?`, yweetID)
	if err != nil {
		return nil, err
	}

	newYweet, err := yweet.ReconstructYweet(uint64(yweetID), insertedYweetDTO.UserID, insertedYweetDTO.Content, insertedYweetDTO.CreatedAt)
	if err != nil {
		return nil, err
	}

	return newYweet, nil
}
