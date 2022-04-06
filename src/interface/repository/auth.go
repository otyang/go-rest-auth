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
	FetchAuth(ctx context.Context, UserID string) (string, error)

	ValidateAccessToken(ctx context.Context, usrID string, sessionID string) error
}

func NewAuthRepository(db *bun.DB, rdb *redis.Client) AuthRepository {
	return &authRepository{db, rdb}
}

func (ar *authRepository) StoreTokenPair(ctx context.Context, td *model.TokenDetails, sessionID string) error {

	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now().UTC()

	sessionIdRt := sessionID + model.PostfixRefreshToken

	// set information in Redis, where the key is the user ID
	if _, err := ar.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, sessionID, "at_id", td.AtID)
		rdb.HSet(ctx, sessionID, "rt_id", td.RtID)
		return nil
	}); err != nil {
		return err
	}
	// set lifetime for redis key
	ar.rdb.Expire(ctx, sessionID, at.Sub(now))

	// set information in Redis, where the key is the user ID
	if _, err := ar.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, sessionIdRt, "rt_id", td.RtID)
		rdb.HSet(ctx, sessionIdRt, "at_id", td.AtID)
		return nil
	}); err != nil {
		return err
	}
	// set lifetime for redis key
	ar.rdb.Expire(ctx, sessionIdRt, rt.Sub(now))

	return nil
}

func (ar *authRepository) StoreAccessToken(ctx context.Context, atd *model.AccessTokenDetails, sessionID string) error {

	twoFactorAuthToken := time.Unix(atd.AtExpires, 0)
	now := time.Now().UTC()

	// set information in Redis, where the key is the user ID
	if _, err := ar.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, sessionID, "at_id", atd.AtID)
		return nil
	}); err != nil {
		return err
	}
	// set lifetime for redis key
	ar.rdb.Expire(ctx, sessionID, twoFactorAuthToken.Sub(now))

	return nil
}

func (ar *authRepository) FetchAuth(ctx context.Context, sessionID string) (string, error) {

	redisRtKey := sessionID + model.PostfixRefreshToken
	if ar.rdb.Exists(ctx, redisRtKey).Val() != 1 {
		return "", errors.New("at token dont found")
	}

	var usrInfo model.UserRedisSessionData
	err := ar.rdb.HMGet(ctx, redisRtKey, "rt_id").Scan(&usrInfo)
	if err != nil {
		return "", err
	}

	if ar.rdb.Exists(ctx, sessionID).Val() == 1 {
		ar.rdb.Del(ctx, sessionID)
	}

	ar.rdb.Del(ctx, redisRtKey)

	return usrInfo.RtID, nil
}

func (ar *authRepository) ValidateAccessToken(ctx context.Context, atID string, sessionID string) error {

	var usrInfo model.UserRedisSessionData
	err := ar.rdb.HMGet(ctx, sessionID, "at_id").Scan(&usrInfo)
	if err != nil {
		return err
	}

	if usrInfo.AtID != atID {
		return errors.New("unauthorized")
	}

	return nil
}
