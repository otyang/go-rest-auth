package repository

import (
	"auth-project/src/domain/model"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type twoFactorAuthRepository struct {
	db *bun.DB
}

type TwoFactorAuthRepository interface {
	SetUpTwoFactorAuthByUserID(ctx context.Context, twoFactorAuthSetUpReq *model.TwoFactorAuthSetUpReq, usrID string) error
	DeleteTwoFactorAuthByUserID(ctx context.Context, twoFactorAuthType, usrID string) error
}

func NewTwoFactorAuthRepository(db *bun.DB) TwoFactorAuthRepository {
	return &twoFactorAuthRepository{db}
}

func (tr *twoFactorAuthRepository) SetUpTwoFactorAuthByUserID(ctx context.Context,
	twoFactorAuthSetUpReq *model.TwoFactorAuthSetUpReq, usrID string) error {

	user := &model.User{ID: usrID}
	err := tr.db.NewSelect().Model(user).
		WherePK().
		Scan(ctx)
	if err != nil {
		return err
	}

	switch twoFactorAuthSetUpReq.Code2faType {
	case model.TokenTypeGoogle:
		if user.IsEmailVerified || user.IsPhoneVerified {
			return errors.New("you can only have one way of two-factor authentication")
		}

		if user.GoogleSecret != "" {
			return errors.New("two-factor google authentication already exist")
		}

		user.GoogleSecret = twoFactorAuthSetUpReq.Secret
		user.IsGoogleVerified = true
		_, err = tr.db.NewUpdate().Model(user).
			WherePK().
			Exec(ctx)
		if err != nil {
			return err
		}

	case model.TokenTypePhone:

		if user.IsEmailVerified || user.GoogleSecret != "" {
			return errors.New("you can only have one way of two-factor authentication")
		}

		if user.IsPhoneVerified {
			return errors.New(" two-factor phone authentication already exist")
		}

		user.IsPhoneVerified = true
		_, err = tr.db.NewUpdate().Model(user).
			WherePK().
			Exec(ctx)
		if err != nil {
			return err
		}

	case model.TokenTypeEmail:

		if user.IsPhoneVerified || user.GoogleSecret != "" {
			return errors.New("you can only have one way of two-factor authentication")
		}

		if user.IsEmailVerified {
			return errors.New("two-factor email authentication already exist")
		}

		user.IsEmailVerified = true
		_, err = tr.db.NewUpdate().Model(user).
			WherePK().
			Exec(ctx)
		if err != nil {
			return err
		}

	default:
		return fiber.NewError(fiber.StatusInternalServerError, "invalid two factor auth type")
	}

	return nil
}

func (tr *twoFactorAuthRepository) DeleteTwoFactorAuthByUserID(ctx context.Context, twoFactorAuthType,
	usrID string) error {

	user := &model.User{ID: usrID}

	switch twoFactorAuthType {
	case model.TokenTypeGoogle:
		_, err := tr.db.NewUpdate().Model(user).
			WherePK().
			Set("google_secret = NULL").
			Set("is_google_verified = FALSE").
			Exec(ctx)
		if err != nil {
			return err
		}
	case model.TokenTypePhone:
		_, err := tr.db.NewUpdate().Model(user).
			WherePK().
			Set("is_phone_verified = FALSE").
			Exec(ctx)
		if err != nil {
			return err
		}
	case model.TokenTypeEmail:
		_, err := tr.db.NewUpdate().Model(user).
			WherePK().
			Set("is_email_verified = FALSE").
			Exec(ctx)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid two factor auth type")
	}

	return nil
}
