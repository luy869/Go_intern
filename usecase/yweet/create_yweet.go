package yweet

import (
	"context"
	"yatter-backend-go/app/domain/object/user"
	"yatter-backend-go/app/domain/object/yweet"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/usecase/transactor"
	"yatter-backend-go/pkg/errors"
)

type CreateYweetUseCase interface {
	Create(ctx context.Context, username, content string) (*yweet.Yweet, *user.User, error)
}

type createYweetResult struct {
	InsertedYweet *yweet.Yweet
	AuthUser      *user.User
}

var _ CreateYweetUseCase = (*CreateYweetUseCaseImpl)(nil)

type CreateYweetUseCaseImpl struct {
	userRepo   repository.User
	yweetRepo  repository.Yweet
	transactor transactor.Transactor
}

func NewCreateYweetUseCase(
	userRepo repository.User,
	yweetRepo repository.Yweet,
	transactor transactor.Transactor,
) *CreateYweetUseCaseImpl {
	return &CreateYweetUseCaseImpl{
		userRepo:   userRepo,
		yweetRepo:  yweetRepo,
		transactor: transactor,
	}
}

func (uc *CreateYweetUseCaseImpl) Create(ctx context.Context, username, content string) (*yweet.Yweet, *user.User, error) {
	result, err := uc.transactor.TransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		authUser, err := uc.userRepo.FindByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if authUser == nil {
			return nil, errors.ErrUnauthorized.WithDevMessage("ユーザーが存在しません")
		}

		pendingYweet, err := yweet.NewPendingYweet(authUser.ID(), content)
		if err != nil {
			return nil, err
		}

		insertedYweet, err := uc.yweetRepo.Insert(ctx, pendingYweet)
		if err != nil {
			return nil, err
		}

		return createYweetResult{
			InsertedYweet: insertedYweet,
			AuthUser:      authUser,
		}, nil
	})

	if err != nil {
		return nil, nil, err
	}

	r, ok := result.(createYweetResult)
	if !ok {
		return nil, nil, errors.ErrInternal.WithDevMessage("failed to cast result")
	}

	return r.InsertedYweet, r.AuthUser, nil
}
