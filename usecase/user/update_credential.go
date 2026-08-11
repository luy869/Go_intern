package user

import (
	"context"
	"yatter-backend-go/app/domain/object/user"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/usecase/transactor"
	"yatter-backend-go/pkg/errors"
)

type UpdateCredentialUseCase interface {
	UpdateCredential(ctx context.Context, username string, newNote string, newDisplayname string) (*user.User, error)
}

var _ UpdateCredentialUseCase = (*UpdateCredentialUseCaseImpl)(nil)

type UpdateCredentialUseCaseImpl struct {
	userRepo   repository.User
	transactor transactor.Transactor
}

func NewUpdateCredentialUseCase(
	userRepo repository.User,
	transactor transactor.Transactor,
) *UpdateCredentialUseCaseImpl {
	return &UpdateCredentialUseCaseImpl{
		userRepo:   userRepo,
		transactor: transactor,
	}
}

func (uc *UpdateCredentialUseCaseImpl) UpdateCredential(ctx context.Context, username string, newNote string, newDisplayname string) (*user.User, error) {
	result, err := uc.transactor.TransactionWithValue(ctx, func(ctx context.Context) (any, error) {
		authUser, err := uc.userRepo.FindByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if authUser == nil {
			return nil, errors.ErrUnauthorized.WithDevMessage("ユーザーが存在しません")
		}
		err = authUser.SetDisplayName(newDisplayname)
		if err != nil {
			return nil, err
		}
		err = authUser.SetNote(newNote)
		if err != nil {
			return nil, err
		}

		updatedUser, err := uc.userRepo.Update(ctx, authUser)
		if err != nil {
			return nil, err
		}
		return updatedUser, nil
	})
	if err != nil {
		return nil, err
	}

	user, ok := result.(*user.User)
	if !ok {
		return nil, errors.ErrInternal.WithDevMessage("failed to cast result to user.User")
	}

	return user, nil
}
