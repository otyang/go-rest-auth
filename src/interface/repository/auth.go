package repository

import (
	"auth-project/src/domain/model"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/uptrace/bun"
	"time"
)

type authRepository struct {
	db  *bun.DB
	rdb *redis.Client
}

type AuthRepository interface {
	StoreTokenPair(ctx context.Context, td *model.TokenDetails, sessionID string) error
	StoreAccessToken(ctx context.Context, atd *model.AccessTokenDetails, sessionID string) error
	FetchAuth(ctx context.Context, sessionID string) (string, error)

	ValidateAccessToken(ctx context.Context, atID string, sessionID string) error
}

func NewAuthRepository(db *bun.DB, rdb *redis.Client) AuthRepository {
	return &authRepository{db, rdb}
}

func (ar *authRepository) StoreTokenPair(ctx context.Context, td *model.TokenDetails, sessionID string) error {

	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now().UTC()

	sessionIdRt := sessionID + model.PostfixRefreshToken

	// set information in Redis, where the key is the session ID
	ar.rdb.Set(ctx, sessionID, td.AtID, at.Sub(now))
	ar.rdb.Set(ctx, sessionIdRt, td.RtID, rt.Sub(now))

	return nil
}

func (ar *authRepository) StoreAccessToken(ctx context.Context, atd *model.AccessTokenDetails, sessionID string) error {

	twoFactorAuthToken := time.Unix(atd.AtExpires, 0)
	now := time.Now().UTC()

	// set information in Redis, where the key is the session ID
	ar.rdb.Set(ctx, sessionID, atd.AtID, twoFactorAuthToken.Sub(now))

	return nil
}

func (ar *authRepository) FetchAuth(ctx context.Context, sessionID string) (string, error) {

	redisRtKey := sessionID + model.PostfixRefreshToken
	if ar.rdb.Exists(ctx, redisRtKey).Val() != 1 {
		return "", errors.New("at token dont found")
	}

	val, err := ar.rdb.Get(ctx, redisRtKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	if ar.rdb.Exists(ctx, sessionID).Val() == 1 {
		ar.rdb.Del(ctx, sessionID)
	}

	ar.rdb.Del(ctx, redisRtKey)

	return val, nil
}

func (ar *authRepository) ValidateAccessToken(ctx context.Context, atID string, sessionID string) error {

	val, err := ar.rdb.Get(ctx, sessionID).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if val != atID {
		return errors.New("unauthorized")
	}

	return nil
}
