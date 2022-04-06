package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis/v8"

	"auth-project/src/domain/model"
	"github.com/go-pg/pg/v10"

	qrcode "github.com/skip2/go-qrcode"
)

type authRepository struct {
	db  *bun.DB
	rdb *redis.Client
}

type AuthRepository interface {
	GetUserByEmailAndPassword(ctx context.Context, email, psw string) (*model.User, error)
	StoreAuth(ctx context.Context, usrInfo *model.UserRedisSessionData, td *model.TokenDetails) error
	GenerateQrCode(ctx context.Context) ([]byte, string, error)
	FetchAuth(ctx context.Context, UserID string) (string, error)
	ValidateAccessTokenID(ctx context.Context, usrID string, atID string) error
}

func NewAuthRepository(db *bun.DB, rdb *redis.Client) AuthRepository {
	return &authRepository{db, rdb}
}

// The GetUserByEmailAndPassword method retrieves User entity
// of domain layer from database by Email and Password.
func (ar *authRepository) GetUserByEmailAndPassword(ctx context.Context,
	email, psw string) (*model.User, error) {

	email = strings.ToLower(email)
	usr := &model.User{}
	err := ar.db.NewSelect().Model(usr).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, errors.New("user not identified")
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(psw))
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (ar *authRepository) StoreAuth(ctx context.Context, usrInfo *model.UserRedisSessionData, td *model.TokenDetails) error {

	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now().UTC()

	// set information in Redis, where the key is the user ID
	if _, err := ar.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, td.AtUuid, "user_id", usrInfo.UserID)
		rdb.HSet(ctx, td.AtUuid, "rt_id", td.RtUuid)
		return nil
	}); err != nil {
		return err
	}
	// set lifetime for redis key
	ar.rdb.Expire(ctx, td.AtUuid, at.Sub(now))

	// set information in Redis, where the key is the user ID
	if _, err := ar.rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, td.RtUuid, "user_id", usrInfo.UserID)
		rdb.HSet(ctx, td.RtUuid, "at_id", td.AtUuid)
		return nil
	}); err != nil {
		return err
	}
	// set lifetime for redis key
	ar.rdb.Expire(ctx, td.RtUuid, rt.Sub(now))

	if ar.rdb.Exists(ctx, usrInfo.UserID).Val() == 1 {
		userSessionsString := ar.rdb.Get(ctx, usrInfo.UserID).Val()
		if userSessionsString == "" {
			return errors.New("userSessionsArray is missing")
		}
		userSessionsString += "," + td.RtUuid
		ar.rdb.Set(ctx, usrInfo.UserID, userSessionsString, rt.Sub(now))
	} else {
		ar.rdb.Set(ctx, usrInfo.UserID, td.RtUuid, rt.Sub(now))
	}

	return nil
}
func (ar *authRepository) GenerateQrCode(ctx context.Context) ([]byte, string, error) {
	token := uuid.New()
	pngByte, err := qrcode.Encode("http://"+
		viper.GetString("http.host")+
		viper.GetString("http.port")+
		"/api/v1/auth/authenticate/qr-code/"+
		token.String(), qrcode.Medium, 256)
	if err != nil {
		return nil, "", err
	}

	return pngByte, token.String(), nil
}

func (ar *authRepository) FetchAuth(ctx context.Context, rtId string) (string, error) {

	var usrInfo model.UserRedisSessionData
	err := ar.rdb.HMGet(ctx, rtId, "user_id", "at_id").Scan(&usrInfo)
	if err != nil {
		return "", err
	}

	if ar.rdb.Exists(ctx, usrInfo.AtID).Val() != 1 {
		return "", errors.New("at token dont found")
	}

	ar.rdb.Del(ctx, usrInfo.AtID)

	return usrInfo.UserID, nil
}

func (ar *authRepository) ValidateAccessTokenID(ctx context.Context, usrID string, atID string) error {

	var usrInfo model.UserRedisSessionData
	err := ar.rdb.HMGet(ctx, atID, "user_id").Scan(&usrInfo)
	if err != nil {
		return err
	}

	if usrInfo.UserID != usrID {
		return errors.New("unauthorized")
	}

	return nil
}
