package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis/v8"

	"auth-project/src/domain/model"
)

type userRepository struct {
	db  *bun.DB
	rdb *redis.Client
}

type UserRepository interface {
	CreateUser(ctx context.Context, usrCrtReq *model.UserCreateReq) (string, error)
	ChangeUserPassword(ctx context.Context, reqData *model.UserChangePasswordReq) error
	UpdateUserByID(ctx context.Context, updReq *model.UserUpdateReq, userID uuid.UUID) (*model.User, error)
	UpdatePinById(ctx context.Context, updReq *model.UserPinUpdateReq, userID uuid.UUID) error
	DeletePinByID(ctx context.Context, updReq *model.UserPinDeleteReq, userID uuid.UUID) error
	SignOut(ctx context.Context, atID string) error
	SignOutAll(ctx context.Context, usrID string) error
}

func NewUserRepository(db *bun.DB, rdb *redis.Client) UserRepository {
	return &userRepository{db, rdb}
}

// The CreateUser method receives a UserCreateReq json object and stores
// User entity in the database.
func (ur *userRepository) CreateUser(ctx context.Context, usrCrtReq *model.UserCreateReq) (string, error) {

	email := strings.ToLower(usrCrtReq.Email)
	exists, err := ur.db.NewSelect().Model((*model.User)(nil)).
		Where("email = ?", email).
		Exists(ctx)
	if err != nil {
		return "", err
	}
	if exists {
		err = errors.New("a user with this email address already exists")
		return "", err
	}

	// Use GenerateFromPassword to hash & salt password.
	hash, err := bcrypt.GenerateFromPassword([]byte(usrCrtReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// New obj:
	usr := &model.User{
		Email:    email,
		Password: string(hash),
	}
	_, err = ur.db.NewInsert().Model(usr).
		Exec(ctx)
	if err != nil {
		return "", err
	}
	return usr.ID.String(), nil
}

// The ChangeUserPassword method retrieves User entity
// of domain layer from database and change password.
func (ur *userRepository) ChangeUserPassword(ctx context.Context, reqData *model.UserChangePasswordReq) error {

	usr := &model.User{ID: reqData.ID}
	err := ur.db.NewSelect().Model(usr).
		WherePK().
		Scan(ctx)
	if err != nil {
		err = errors.New("user not identified")
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(reqData.OldPassword))
	if err != nil {
		err = errors.New("incorrect current password")
		return err
	}

	// Use GenerateFromPassword to hash & salt password.
	hash, err := bcrypt.GenerateFromPassword([]byte(reqData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usr.Password = string(hash)
	_, err = ur.db.NewUpdate().Model(usr).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// The UpdateUserByID method update User entity
// of domain layer from database.
func (ur *userRepository) UpdateUserByID(ctx context.Context, updReq *model.UserUpdateReq,
	userID uuid.UUID) (*model.User, error) {

	user := &model.User{ID: userID}
	err := ur.db.NewSelect().Model(user).
		WherePK().
		Scan(ctx)

	user.FullName = updReq.FullName
	user.FullName = updReq.Email
	user.FullName = updReq.Phone

	_, err = ur.db.NewUpdate().Model(user).
		WherePK().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// The UpdatePinById method update pin in User entity
// of domain layer from database.
func (ur *userRepository) UpdatePinById(ctx context.Context, updReq *model.UserPinUpdateReq, userID uuid.UUID) error {

	user := &model.User{ID: userID}
	err := ur.db.NewSelect().
		Model(user).
		WherePK().
		Scan(ctx)

	if user.Pin != "" {
		err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(updReq.OldPin))
		if err != nil {
			return errors.New("incorrect email")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(updReq.NewPin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Pin = string(hash)
	_, err = ur.db.NewUpdate().Model(user).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// The DeletePinByID method delete pin in User entity
// of domain layer from database.
func (ur *userRepository) DeletePinByID(ctx context.Context, updReq *model.UserPinDeleteReq, userID uuid.UUID) error {

	user := &model.User{ID: userID}
	err := ur.db.NewSelect().Model(user).
		WherePK().
		Scan(ctx)

	err = bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(updReq.Pin))
	if err != nil {
		return errors.New("incorrect email")
	}

	_, err = ur.db.NewUpdate().Model(user).
		WherePK().
		Set("pin = NULL").
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// SignOut clear redis key, and check exist
func (ur *userRepository) SignOut(ctx context.Context, atID string) error {

	var usrInfo model.UserRedisSessionData
	err := ur.rdb.HMGet(ctx, atID, "user_id", "rt_id").Scan(&usrInfo)
	if err != nil {
		return err
	}

	ur.rdb.Del(ctx, atID)
	ur.rdb.Del(ctx, usrInfo.RtID)

	if ur.rdb.Exists(ctx, atID).Val() == 1 ||
		ur.rdb.Exists(ctx, usrInfo.RtID).Val() == 1 {
		return errors.New("redis key deleting error")
	}

	return nil
}

// SignOutAll clear redis key, and check exist
func (ur *userRepository) SignOutAll(ctx context.Context, usrID string) error {

	userSessionsString := ur.rdb.Get(ctx, usrID).Val()
	if userSessionsString == "" {
		return errors.New("userSessionsString is missing")
	}

	for _, id := range strings.Split(userSessionsString, ",") {
		var usrInfo model.UserRedisSessionData
		err := ur.rdb.HMGet(ctx, id, "user_id", "at_id").Scan(&usrInfo)
		if err != nil {
			return err
		}

		ur.rdb.Del(ctx, id)
		ur.rdb.Del(ctx, usrInfo.AtID)

		if ur.rdb.Exists(ctx, id).Val() == 1 ||
			ur.rdb.Exists(ctx, usrInfo.AtID).Val() == 1 {
			return errors.New("redis key deleting error")
		}
	}

	return nil
}
