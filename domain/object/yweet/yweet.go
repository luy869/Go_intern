package yweet

import (
	"time"
	"yatter-backend-go/pkg/errors"
)

type Yweet struct {
	id        uint64
	userID    uint64
	content   string
	createdAt time.Time
}

func ReconstructYweet(id, userID uint64, content string, createdAt time.Time) (*Yweet, error) {
	yweet := &Yweet{}

	if err := yweet.SetID(id); err != nil {
		return nil, err
	}

	if err := yweet.SetUserID(userID); err != nil {
		return nil, err
	}

	if err := yweet.SetContent(content); err != nil {
		return nil, err
	}

	if err := yweet.SetCreatedAt(createdAt); err != nil {
		return nil, err
	}

	return yweet, nil
}

func (y *Yweet) SetID(id uint64) error {
	if id < 1 {
		return errors.ErrInternal.WithDevMessage("id must be more than 0")
	}

	y.id = id
	return nil
}

func (y *Yweet) SetUserID(userID uint64) error {
	if userID < 1 {
		return errors.ErrInternal.WithDevMessage("IDが不正です")
	}

	y.userID = userID
	return nil
}

func (y *Yweet) SetContent(content string) error {
	if content == "" {
		return errors.ErrInternal.WithDevMessage("本文は必須です")
	}
	y.content = content
	return nil
}

func (y *Yweet) SetCreatedAt(createdAt time.Time) error {
	yatterLaunchedAt := time.Date(2025, 4, 1, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))
	if !createdAt.After(yatterLaunchedAt) {
		return errors.ErrInternal.WithDevMessage("createdAt must be after yatter launched")
	}

	if createdAt.After(time.Now()) {
		return errors.ErrInternal.WithDevMessage("createdAt must not be in the future")
	}

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	jstCreatedAt := createdAt.In(jst)
	if !createdAt.Equal(jstCreatedAt) {
		return errors.ErrInternal.WithDevMessage("createdAt must be in JST")
	}

	y.createdAt = createdAt
	return nil
}

func (y *Yweet) ID() uint64 {
	return y.id
}

func (y *Yweet) UserID() uint64 {
	return y.userID
}

func (y *Yweet) Content() string {
	return y.content
}

func (y *Yweet) CreatedAt() time.Time {
	return y.createdAt
}
