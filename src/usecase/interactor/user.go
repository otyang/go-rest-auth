package interactor

import (
	"context"
	"github.com/google/uuid"

	"auth-project/src/domain/model"
	"auth-project/src/usecase/presenter"
	"auth-project/src/usecase/repository"
)

type userInteractor struct {
	UserRepository repository.UserRepository
	UserPresenter  presenter.UserPresenter
}

type UserInteractor interface {
	CreateUser(ctx context.Context, usrCrtReq *model.UserCreateReq) (string, error)
	ChangeUserPassword(ctx context.Context, reqData *model.UserChangePasswordReq) error
	UpdateUserByID(ctx context.Context, updReq *model.UserUpdateReq, userID uuid.UUID) (*model.UserUpdResp, error)
	UpdatePinById(ctx context.Context, updReq *model.UserPinUpdateReq, userID uuid.UUID) error
	DeletePinByID(ctx context.Context, updReq *model.UserPinDeleteReq, userID uuid.UUID) error
	SignOut(ctx context.Context, atID string) error
	SignOutAll(ctx context.Context, usrID string) error
}

func NewUserInteractor(
	r repository.UserRepository, p presenter.UserPresenter) UserInteractor {
	return &userInteractor{r, p}
}

func (ui *userInteractor) CreateUser(ctx context.Context, usrCrtReq *model.UserCreateReq) (string, error) {

	id, err := ui.UserRepository.CreateUser(ctx, usrCrtReq)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (ui *userInteractor) ChangeUserPassword(ctx context.Context, reqData *model.UserChangePasswordReq) error {

	err := ui.UserRepository.ChangeUserPassword(ctx, reqData)
	if err != nil {
		return err
	}
	return nil
}

func (ui *userInteractor) UpdatePinById(ctx context.Context, updReq *model.UserPinUpdateReq, userID uuid.UUID) error {

	err := ui.UserRepository.UpdatePinById(ctx, updReq, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ui *userInteractor) DeletePinByID(ctx context.Context, updReq *model.UserPinDeleteReq, userID uuid.UUID) error {

	err := ui.UserRepository.DeletePinByID(ctx, updReq, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ui *userInteractor) UpdateUserByID(ctx context.Context, updReq *model.UserUpdateReq,
	userID uuid.UUID) (*model.UserUpdResp, error) {

	user, err := ui.UserRepository.UpdateUserByID(ctx, updReq, userID)
	if err != nil {
		return nil, err
	}
	return ui.UserPresenter.UpdateUserByIDResp(user), nil
}

func (ui *userInteractor) SignOut(ctx context.Context, atID string) error {

	err := ui.UserRepository.SignOut(ctx, atID)
	if err != nil {
		return err
	}
	return nil
}

func (ui *userInteractor) SignOutAll(ctx context.Context, usrID string) error {

	err := ui.UserRepository.SignOutAll(ctx, usrID)
	if err != nil {
		return err
	}
	return nil
}
