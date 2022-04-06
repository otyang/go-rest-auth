package repository

import (
	"auth-project/tools"
	"context"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/uptrace/bun"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis/v8"

	"auth-project/src/domain/model"
)

type userRepository struct {
	db  *bun.DB
	rdb *redis.Client
}

type UserRepository interface {
	IsExitsUserByEmail(ctx context.Context, email string) (bool, error)
	IsExitsUserByIDAndPassword(ctx context.Context, userID, password string) (bool, error)

	GetUserByLoginAndPassword(ctx context.Context, login, psw string) (*model.User, error)
	GetUserByEmailOrPhone(ctx context.Context, emailOrPhone string) (*model.User, error)
	GetUserByID(ctx context.Context, usrID string) (*model.User, error)

	CreatUnActivateUserIfNotExistByLogin(ctx context.Context, data *model.SignUpSend2faCodeReq) error
	SignUpActivateUser(ctx context.Context, signUpReq *model.SignUpActivateUserData) error

	ResetUserPassword(ctx context.Context, newPassword, usrID string) error

	ChangeUserPasswordByID(ctx context.Context, data *model.UserChangePasswordData, usrID string) error

	UpdateUserInfoByID(ctx context.Context, updReq *model.UserUpdateInfoData, userID string) (*model.User, error)
	UpdateUserEmailByID(ctx context.Context, email string, userID string) (*model.User, error)
	UpdateUserPhoneByID(ctx context.Context, phone string, userID string) (*model.User, error)

	SignOut(ctx context.Context, atID string) error
	SignOutAll(ctx context.Context, usrID string) error
}

func NewUserRepository(db *bun.DB, rdb *redis.Client) UserRepository {
	return &userRepository{db, rdb}
}

func (ur *userRepository) IsExitsUserByEmail(ctx context.Context, email string) (bool, error) {

	exists, err := ur.db.NewSelect().Model((*model.User)(nil)).
		Where("email = ? ", email).
		Exists(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (ur *userRepository) IsExitsUserByIDAndPassword(ctx context.Context, userID string, password string) (bool, error) {

	usr := &model.User{ID: userID}
	err := ur.db.NewSelect().Model(usr).
		WherePK().
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("incorrect id")
		}

		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		return false, errors.New("incorrect password")
	}

	return true, nil
}

// The GetUserByLoginAndPassword method retrieves User entity
// of domain layer from database by Email and Password.
func (ur *userRepository) GetUserByLoginAndPassword(ctx context.Context,
	login, psw string) (*model.User, error) {

	login = strings.ToLower(login)
	usr := &model.User{}
	err := ur.db.NewSelect().Model(usr).
		Where("is_active = true").
		Where("email = ? OR phone = ? OR id = ?", login, login, login).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("incorrect login or password")
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(psw))
	if err != nil {
		return nil, errors.New("incorrect login or password")
	}

	return usr, nil
}

func (ur *userRepository) GetUserByEmailOrPhone(ctx context.Context, emailOrPhone string) (*model.User, error) {

	emailOrPhone = strings.ToLower(emailOrPhone)
	usr := &model.User{}
	err := ur.db.NewSelect().Model(usr).
		Where("is_active = true").
		Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (ur *userRepository) GetUserByID(ctx context.Context, usrID string) (*model.User, error) {

	usr := &model.User{ID: usrID}
	err := ur.db.NewSelect().Model(usr).
		Relation("ReferralUser").
		WherePK().
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not identified")
		}

		return nil, err
	}

	return usr, nil
}

func (ur *userRepository) CreatUnActivateUserIfNotExistByLogin(ctx context.Context, data *model.SignUpSend2faCodeReq) error {

	var usr model.User
	err := ur.db.NewSelect().Model(&usr).
		Where("email = ? OR phone = ?", data.Login, data.Login).
		Scan(ctx)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			id, err := gonanoid.New()
			if err != nil {
				return err
			}
			user := model.User{
				ID:           id,
				ReferralLink: tools.GenerateLink(),
				IsActive:     false,
			}
			switch data.LoginType {
			case model.TokenTypePhone:
				user.Phone = data.Login
			case model.TokenTypeEmail:
				user.Email = data.Login
			default:
				return fiber.NewError(fiber.StatusBadRequest, "invalid login type")
			}

			_, err = ur.db.NewInsert().Model(&user).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
	}

	if usr.IsActive {
		return errors.New("user already exist and activate")
	}

	return nil
}

func (ur *userRepository) SignUpActivateUser(ctx context.Context, signUpReq *model.SignUpActivateUserData) error {

	var usr model.User
	err := ur.db.NewSelect().Model(&usr).
		Where("email = ? OR phone = ?", signUpReq.Login, signUpReq.Login).
		Scan(ctx)
	if err != nil {
		return err
	}

	if usr.IsActive {
		return errors.New("user already activate")
	}

	// Use GenerateFromPassword to hash & salt password.
	hash, err := bcrypt.GenerateFromPassword([]byte(signUpReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usr.Password = string(hash)
	usr.IsActive = true
	usr.UpdatedAt = time.Now().UTC()

	if signUpReq.Referral != "" {
		exists, err := ur.db.NewSelect().Model((*model.User)(nil)).
			Where("referral_link = ? ", signUpReq.Referral).
			Exists(ctx)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("referral link invalid")
		}

		usr.Referral = signUpReq.Referral
	}

	_, err = ur.db.NewUpdate().Model(&usr).
		WherePK().
		OmitZero().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) ResetUserPassword(ctx context.Context, newPassword, usrID string) error {

	usr := &model.User{ID: usrID}
	err := ur.db.NewSelect().Model(usr).
		WherePK().
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("user not identified")
		}
		return err
	}

	// Use GenerateFromPassword to hash & salt password.
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usr.Password = string(hash)
	usr.UpdatedAt = time.Now().UTC()
	_, err = ur.db.NewUpdate().Model(usr).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// The ChangeUserPasswordByID method retrieves User entity
// of domain layer from database and change password.
func (ur *userRepository) ChangeUserPasswordByID(ctx context.Context, data *model.UserChangePasswordData,
	usrID string) error {

	usr := &model.User{ID: usrID}
	err := ur.db.NewSelect().Model(usr).
		WherePK().
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("user not identified")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(data.OldPassword))
	if err != nil {
		err = errors.New("incorrect current password")
		return err
	}

	// Use GenerateFromPassword to hash & salt password.
	hash, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	usr.Password = string(hash)
	usr.UpdatedAt = time.Now().UTC()
	_, err = ur.db.NewUpdate().Model(usr).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUserInfoByID(ctx context.Context, updReq *model.UserUpdateInfoData,
	userID string) (*model.User, error) {

	user := &model.User{
		ID:       userID,
		FullName: updReq.FullName,
	}
	_, err := ur.db.NewUpdate().Model(user).
		Column("full_name").
		WherePK().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) UpdateUserPhoneByID(ctx context.Context, phone string,
	userID string) (*model.User, error) {

	user := &model.User{
		ID:    userID,
		Phone: phone,
	}
	_, err := ur.db.NewUpdate().Model(user).
		Column("phone").
		WherePK().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepository) UpdateUserEmailByID(ctx context.Context, email string,
	userID string) (*model.User, error) {

	user := &model.User{
		ID:    userID,
		Email: email,
	}
	_, err := ur.db.NewUpdate().Model(user).
		Column("email").
		WherePK().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// SignOut clear redis key, and check exist
func (ur *userRepository) SignOut(ctx context.Context, sessionID string) error {

	var usrInfo model.UserRedisSessionData
	err := ur.rdb.HMGet(ctx, sessionID, "user_id", "rt_id").Scan(&usrInfo)
	if err != nil {
		return err
	}

	ur.rdb.Del(ctx, sessionID)
	ur.rdb.Del(ctx, sessionID+model.PostfixRefreshToken)

	if ur.rdb.Exists(ctx, sessionID).Val() == 1 ||
		ur.rdb.Exists(ctx, sessionID+model.PostfixRefreshToken).Val() == 1 {
		return errors.New("redis key deleting error")
	}

	_, err = ur.db.NewUpdate().Model((*model.Session)(nil)).
		Where("session_id = ?", sessionID).
		Set("is_logout = TRUE").
		Set("updated_at = ?", time.Now().UTC()).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// SignOutAll clear redis key, and check exist
func (ur *userRepository) SignOutAll(ctx context.Context, usrID string) error {

	var sessions []model.Session
	err := ur.db.NewSelect().Model(&sessions).
		Where("user_id = ?", usrID).
		Where("is_logout = FALSE").
		Where("expires_at > ?", time.Now().UTC()).
		Scan(ctx)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		var usrInfo model.UserRedisSessionData
		err := ur.rdb.HMGet(ctx, session.SessionID+model.PostfixRefreshToken, "user_id", "at_id").Scan(&usrInfo)
		if err != nil {
			return err
		}

		ur.rdb.Del(ctx, session.SessionID+model.PostfixRefreshToken)
		ur.rdb.Del(ctx, session.SessionID)

		if ur.rdb.Exists(ctx, session.SessionID+model.PostfixRefreshToken).Val() == 1 ||
			ur.rdb.Exists(ctx, session.SessionID).Val() == 1 {
			return errors.New("redis key deleting error")
		}
	}

	_, err = ur.db.NewUpdate().Model(&sessions).
		Where("expires_at > ?", time.Now().UTC()).
		Set("is_logout = TRUE").
		Set("updated_at = ?", time.Now().UTC()).
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
