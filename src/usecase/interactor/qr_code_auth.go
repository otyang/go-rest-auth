package interactor

import (
	"auth-project/src/domain/model"
	"auth-project/src/infrastructure/authentication"
	"auth-project/src/usecase/presenter"
	"auth-project/src/usecase/repository"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

type qrCodeAuthInteractor struct {
	AuthRepository       repository.AuthRepository
	SessionRepository    repository.SessionRepository
	UserRepository       repository.UserRepository
	QrCodeAuthRepository repository.QrCodeAuthRepository

	QrCodeAuthPresenter presenter.QrCodeAuthPresenter

	jwtConfigurator *authentication.JwtConfigurator
}

type QrCodeAuthInteractor interface {
	GenerateQrCode(ctx context.Context) ([]byte, string, error)
	CreateAuthToken(ctx context.Context, usrID string) (string, error)
	GetUserByQrCodeAuthToken(ctx context.Context, token string) (*model.User, error)
	GenerateTokenPairByUserID(c *websocket.Conn, usrID string) (*model.TokenDetails, error)
}

func NewQrCodeAuthInteractor(
	ar repository.AuthRepository, sr repository.SessionRepository, ur repository.UserRepository, qr repository.QrCodeAuthRepository, p presenter.QrCodeAuthPresenter, jc *authentication.JwtConfigurator) QrCodeAuthInteractor {
	return &qrCodeAuthInteractor{ar, sr, ur, qr, p, jc}
}

func (qi *qrCodeAuthInteractor) GenerateQrCode(ctx context.Context) ([]byte, string, error) {

	qrCode, tokenID, err := qi.QrCodeAuthRepository.GenerateQrCode(ctx)
	if err != nil {
		return nil, "", err
	}
	return qrCode, tokenID, nil
}

func (qi *qrCodeAuthInteractor) CreateAuthToken(ctx context.Context, usrID string) (string, error) {

	authToken, err := qi.QrCodeAuthRepository.CreateAuthToken(ctx, usrID)
	if err != nil {
		return "", err
	}
	return authToken, nil
}

func (qi *qrCodeAuthInteractor) GetUserByQrCodeAuthToken(ctx context.Context, token string) (*model.User, error) {

	user, err := qi.QrCodeAuthRepository.GetUserByQrCodeAuthToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (qi *qrCodeAuthInteractor) GenerateTokenPairByUserID(c *websocket.Conn, usrID string) (*model.TokenDetails, error) {

	// Make access and refresh tokens to user
	usr, err := qi.UserRepository.GetUserByID(context.Background(), usrID)
	if err != nil {
		return nil, err
	}

	userAgent, ok := c.Locals("User-Agent").(string)
	if !ok {
		return nil, errors.New("err load user agent")
	}

	sessionID, err := gonanoid.New()
	if err != nil {
		return nil, err
	}

	details, err := qi.jwtConfigurator.GenerateTokenPair(usr.ID, sessionID)
	if err != nil {
		return nil, err
	}

	err = qi.AuthRepository.StoreTokenPair(context.Background(), details, sessionID)
	if err != nil {
		return nil, err
	}

	ses := &model.Session{
		SessionID: sessionID,
		UserAgent: userAgent,
		ClientIP:  c.Conn.LocalAddr().String(),
		ExpiresAT: time.Unix(details.RtExpires, 0).UTC(),
		UserID:    usr.ID,
	}

	err = qi.SessionRepository.InsertSession(context.Background(), ses)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return details, nil
}
