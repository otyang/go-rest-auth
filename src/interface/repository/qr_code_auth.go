package repository

import (
	"auth-project/src/domain/model"
	"auth-project/tools"
	"context"
	"database/sql"
	"errors"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"time"
)

type qrCodeAuthRepository struct {
	db *bun.DB
}

type QrCodeAuthRepository interface {
	GenerateQrCode(ctx context.Context) ([]byte, string, error)
	CreateAuthToken(ctx context.Context, usrID string) (string, error)
	GetUserByQrCodeAuthToken(ctx context.Context, token string) (*model.User, error)
}

func NewQrCodeAuthRepository(db *bun.DB) QrCodeAuthRepository {
	return &qrCodeAuthRepository{db}
}

func (qr *qrCodeAuthRepository) GenerateQrCode(ctx context.Context) ([]byte, string, error) {
	token, err := gonanoid.New()
	if err != nil {
		return nil, "", err
	}
	pngByte, err := qrcode.Encode(
		viper.GetString("http_front.host")+
			"/api/v1/auth/authenticate/qr-code/"+
			token, qrcode.Medium, 256)
	if err != nil {
		return nil, "", err
	}

	return pngByte, token, nil
}

func (qr *qrCodeAuthRepository) CreateAuthToken(ctx context.Context, usrID string) (string, error) {
	code := tools.RandStr(255, "alphanum")

	// New obj:
	token := &model.Token{
		Target:    usrID,
		Value:     code,
		Reason:    model.TokenReasonAuthByQrCode,
		ExpiresAT: tools.AddTimeToCurrentDate(viper.GetDuration("qr_code.token_min_lifetime")),
	}

	var err error
	token.ID, err = gonanoid.New()
	if err != nil {
		return "", err
	}

	_, err = qr.db.NewInsert().Model(token).
		Exec(ctx)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (qr *qrCodeAuthRepository) GetUserByQrCodeAuthToken(ctx context.Context, authToken string) (*model.User, error) {
	var token model.Token
	err := qr.db.NewSelect().Model(&token).
		Where("value = ? AND reason = ? AND expires_at > ? AND is_used = FALSE",
			authToken, model.TokenReasonAuthByQrCode, time.Now().UTC()).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("authToken invalid")
		}
		return nil, err
	}

	usr := &model.User{}
	err = qr.db.NewSelect().Model(usr).
		Where("id = ? ", token.Target).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	_, err = qr.db.NewUpdate().Model(&token).
		WherePK().
		Set("is_used = TRUE").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
