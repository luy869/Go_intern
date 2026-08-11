package yweet

import "yatter-backend-go/pkg/errors"

type PendingYweet struct {
	userID  uint64
	content string
}

func NewPendingYweet(userID uint64, content string) (*PendingYweet, error) {
	yweet := &PendingYweet{}

	if err := yweet.SetUserID(userID); err != nil {
		return nil, err
	}

	if err := yweet.SetContent(content); err != nil {
		return nil, err
	}

	return yweet, nil
}

func (y *PendingYweet) SetUserID(userID uint64) error {
	if userID < 1 {
		return errors.ErrInternal.WithDevMessage("IDが不正です")
	}

	y.userID = userID
	return nil
}

func (y *PendingYweet) SetContent(content string) error {
	if content == "" {
		return errors.ErrBadRequest.WithDevMessage("本文は必須です")
	}
	y.content = content
	return nil
}

func (y *PendingYweet) UserID() uint64 {
	return y.userID
}

func (y *PendingYweet) Content() string {
	return y.content
}
