package interactor

import (
	"context"

	"auth-project/src/domain/model"
	"auth-project/src/usecase/presenter"
	"auth-project/src/usecase/repository"
)

type authInteractor struct {
	AuthRepository repository.AuthRepository
	AuthPresenter  presenter.AuthPresenter
}

type AuthInteractor interface {
	GetUserByEmailAndPassword(ctx context.Context, email, psw string) (*model.User, error)
	StoreAuth(ctx context.Context, usrInfo *model.UserRedisSessionData, td *model.TokenDetails) error
	GenerateQrCode(ctx context.Context) ([]byte, string, error)
	FetchAuth(ctx context.Context, rtId string) (string, error)
	ValidateAccessTokenID(ctx context.Context, usrID string, atID string) error
}

func NewAuthInteractor(
	r repository.AuthRepository, p presenter.AuthPresenter) AuthInteractor {
	return &authInteractor{r, p}
}

func (ai *authInteractor) GetUserByEmailAndPassword(ctx context.Context,
	email, psw string) (*model.User, error) {

	usr, err := ai.AuthRepository.GetUserByEmailAndPassword(ctx, email, psw)
	if err != nil {
		return nil, err
	}
	return ai.AuthPresenter.GetUserByEmailAndPasswordResp(usr), nil
}

func (ai *authInteractor) StoreAuth(ctx context.Context, usrInfo *model.UserRedisSessionData, td *model.TokenDetails) error {

	err := ai.AuthRepository.StoreAuth(ctx, usrInfo, td)
	if err != nil {
		return err
	}
	return err
}
func (ai *authInteractor) GenerateQrCode(ctx context.Context) ([]byte, string, error) {

	qrCode, tokenID, err := ai.AuthRepository.GenerateQrCode(ctx)
	if err != nil {
		return nil, "", err
	}
	return qrCode, tokenID, nil
}

func (ai *authInteractor) FetchAuth(ctx context.Context, rtId string) (string, error) {

	usrID, err := ai.AuthRepository.FetchAuth(ctx, rtId)
	if err != nil {
		return "", err
	}
	return usrID, err
}

func (ai *authInteractor) ValidateAccessTokenID(ctx context.Context, usrID string, atID string) error {

	err := ai.AuthRepository.ValidateAccessTokenID(ctx, usrID, atID)
	if err != nil {
		return err
	}
	return nil
}
