package repository

import (
	"auth-project/src/domain/model"
	"auth-project/src/infrastructure/sending/email"
	"auth-project/src/infrastructure/sending/sms"
	"auth-project/tools"
	"context"
	"database/sql"
	"errors"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"time"
)

type tokenRepository struct {
	db *bun.DB
}

type TokenRepository interface {
	Validate2faCode(ctx context.Context, verifyCodeDate *model.VerifyCodeData) (*model.Token, error)
	Send2faCode(ctx context.Context, sendCodeDate *model.Send2faCodeData) (string, error)
	TokenSetUsed(ctx context.Context, verifyCodeDate *model.VerifyCodeData) error
}

func NewTokenRepository(db *bun.DB) TokenRepository {
	return &tokenRepository{db}
}

func (tr *tokenRepository) Validate2faCode(ctx context.Context, verifyCodeDate *model.VerifyCodeData) (*model.Token, error) {

	if verifyCodeDate.UserID == "" && verifyCodeDate.Target == "" {
		return nil, errors.New("user id and target empty")
	}

	var token model.Token
	query := tr.db.NewSelect().Model(&token).
		Where("is_used = FALSE").
		Where("value = ? ", verifyCodeDate.Code).
		Where("reason = ? ", verifyCodeDate.Reason)
	if verifyCodeDate.Target != "" {
		query = query.Where("target = ? ", verifyCodeDate.Target)
	}
	if verifyCodeDate.CodeType != "" {
		query = query.Where("type = ? ", verifyCodeDate.CodeType)
	}
	if verifyCodeDate.UserID != "" {
		query = query.Where("user_id = ? ", verifyCodeDate.UserID)
	}
	err := query.Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid token")
		}
		return nil, err
	}
	return &token, nil
}

func (tr *tokenRepository) Send2faCode(ctx context.Context, sendOTPDate *model.Send2faCodeData) (string, error) {

	query := tr.db.NewSelect().Model((*model.Token)(nil))
	if sendOTPDate.UserID != "" {
		query = query.Where("user_id = ?", sendOTPDate.UserID)
	} else {
		query = query.Where("target = ?", sendOTPDate.Target)
	}
	exists, err := query.
		Where("is_used = ?", false).
		Where("reason = ?", sendOTPDate.Reason).
		Where("created_at > ?", time.Now().UTC().Add(-viper.GetDuration("2fa.send_timeout"))).
		Exists(ctx)
	if err != nil {
		return "", err
	}

	if exists {
		return "", errors.New(model.TokenTimeSendErr)
	}

	code := tools.RandStr(6, "number")
	// New obj:
	token := &model.Token{
		Target:    sendOTPDate.Target,
		Value:     code,
		Reason:    sendOTPDate.Reason,
		Type:      sendOTPDate.Code2faType,
		ExpiresAT: tools.AddTimeToCurrentDate(viper.GetDuration("2fa.token_min_lifetime")),
	}

	if sendOTPDate.UserID != "" {
		token.UserID = sendOTPDate.UserID
	}

	token.ID, err = gonanoid.New()
	if err != nil {
		return "", err
	}

	_, err = tr.db.NewInsert().Model(token).
		Exec(ctx)
	if err != nil {
		return "", err
	}

	switch sendOTPDate.Code2faType {
	case model.TokenTypePhone:
		err = sms.SendSms(sendOTPDate.Target, sms.CreateSmsBodyVerificationCode(code))
		if err != nil {
			return "", err
		}

	case model.TokenTypeEmail:
		plain, html := email.CreateEmailBodyVerificationCode(code)

		err = email.SendEmail("Verification Code", sendOTPDate.Target, "Verification Code", plain, html)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("invalid token type")
	}

	return sendOTPDate.Target, nil
}

func (tr *tokenRepository) TokenSetUsed(ctx context.Context, verifyCodeDate *model.VerifyCodeData) error {

	if verifyCodeDate.UserID == "" && verifyCodeDate.Target == "" {
		return errors.New("user id and target empty")
	}

	token := model.Token{
		IsUsed: true,
	}
	query := tr.db.NewUpdate().Model(&token).
		Where("is_used = FALSE").
		Where("value = ? ", verifyCodeDate.Code).
		Where("reason = ? ", verifyCodeDate.Reason).
		OmitZero()
	if verifyCodeDate.Target != "" {
		query = query.Where("target = ? ", verifyCodeDate.Target)
	}
	if verifyCodeDate.CodeType != "" {
		query = query.Where("type = ? ", verifyCodeDate.CodeType)
	}
	if verifyCodeDate.UserID != "" {
		query = query.Where("user_id = ? ", verifyCodeDate.UserID)
	}
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("on update, token not found")
	}

	return nil
}
