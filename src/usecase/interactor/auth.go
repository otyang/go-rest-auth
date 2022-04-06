package interactor

import (
	"auth-project/src/domain/model"
	"auth-project/src/infrastructure/authentication"
	"auth-project/src/usecase/presenter"
	"auth-project/src/usecase/repository"
	"context"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

type authInteractor struct {
	AuthRepository    repository.AuthRepository
	SessionRepository repository.SessionRepository
	UserRepository    repository.UserRepository
	TokenRepository   repository.TokenRepository

	AuthPresenter presenter.AuthPresenter

	jwtConfigurator *authentication.JwtConfigurator
}

type AuthInteractor interface {
	Authenticate(ctx context.Context, authReq *model.AuthenticationReq, usrInfo *model.UserSessionData) (map[string]interface{}, error)
	RefreshToken(ctx context.Context, usrInfo *model.UserSessionData, bearerToken string) (map[string]interface{}, error)

	ValidateAccessToken(ctx context.Context, bearerToken string) (*model.AccessClaims, error)
	ValidateTwoFactorAuthToken(ctx context.Context, bearerToken string) (*model.AccessClaims, error)
}

func NewAuthInteractor(
	ar repository.AuthRepository, sr repository.SessionRepository, ur repository.UserRepository, tr repository.TokenRepository, p presenter.AuthPresenter, jc *authentication.JwtConfigurator) AuthInteractor {
	return &authInteractor{ar, sr, ur, tr, p, jc}
}

func (ai *authInteractor) Authenticate(ctx context.Context, authReq *model.AuthenticationReq,
	usrInfo *model.UserSessionData) (map[string]interface{}, error) {

	usr, err := ai.UserRepository.GetUserByLoginAndPassword(ctx,
		authReq.Email, authReq.Password)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	usrInfo.UserID = usr.ID

	sessionID, err := gonanoid.New()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// if two-factor auth token is enabled, then we give a token for two-factor auth
	if usr.IsEmailVerified || usr.IsPhoneVerified || usr.IsGoogleVerified {

		details, err := ai.jwtConfigurator.GenerateTwoFactorAuthToken(usr.ID, sessionID)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = ai.AuthRepository.StoreAccessToken(ctx, details, sessionID)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		switch {
		case usr.IsGoogleVerified:
			return map[string]interface{}{
				"2fa_auth_token": details.AccessToken,
				"2fa_type":       model.TokenTypeGoogle,
			}, nil

		case usr.IsPhoneVerified:
			_, err := ai.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
				UserID:      usr.ID,
				Target:      usr.Phone,
				Code2faType: model.TokenTypePhone,
				Reason:      model.TokenReasonTwoFactorAuth,
			})
			if err != nil {
				if err.Error() != model.TokenTimeSendErr {
					return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
			}

			return map[string]interface{}{
				"2fa_auth_token": details.AccessToken,
				"2fa_type":       model.TokenTypePhone,
				"2fa_target":     usr.Phone,
			}, nil

		case usr.IsEmailVerified:
			_, err := ai.TokenRepository.Send2faCode(ctx, &model.Send2faCodeData{
				UserID:      usr.ID,
				Target:      usr.Email,
				Code2faType: model.TokenTypeEmail,
				Reason:      model.TokenReasonTwoFactorAuth,
			})
			if err != nil {
				if err.Error() != model.TokenTimeSendErr {
					return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
			}

			return map[string]interface{}{
				"2fa_auth_token": details.AccessToken,
				"2fa_type":       model.TokenTypeEmail,
				"2fa_target":     usr.Email,
			}, nil

		default:
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
		}
	}

	details, err := ai.jwtConfigurator.GenerateTokenPair(usr.ID, sessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = ai.AuthRepository.StoreTokenPair(ctx, details, sessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ses := &model.Session{
		SessionID: sessionID,
		UserAgent: usrInfo.UserAgent,
		ClientIP:  usrInfo.ClientIp,
		ExpiresAT: time.Unix(details.RtExpires, 0).UTC(),
		UserID:    usrInfo.UserID,
	}

	err = ai.SessionRepository.InsertSession(ctx, ses)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return map[string]interface{}{
		"access_token":  details.AccessToken,
		"refresh_token": details.RefreshToken,
	}, nil
}

func (ai *authInteractor) RefreshToken(ctx context.Context,
	usrInfo *model.UserSessionData, bearerToken string) (map[string]interface{}, error) {

	claims, err := ai.jwtConfigurator.GetRefreshTokenClaims(bearerToken)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	usrInfo.UserID = claims.UserID

	rtID, err := ai.AuthRepository.FetchAuth(ctx, claims.SessionID)
	if err != nil {
		if err.Error() == "at token dont found" {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if rtID != claims.RtID {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "invalid token")
	}

	// All OK, re-generate the new pair and send to client,
	// we could only generate an access token as well.
	details, err := ai.jwtConfigurator.GenerateTokenPair(claims.UserID, claims.SessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = ai.AuthRepository.StoreTokenPair(ctx, details, claims.SessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ses := &model.Session{
		SessionID: claims.SessionID,
		UserAgent: usrInfo.UserAgent,
		ClientIP:  usrInfo.ClientIp,
		ExpiresAT: time.Unix(details.RtExpires, 0).UTC(),
	}

	err = ai.SessionRepository.UpdateSession(ctx, ses)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return map[string]interface{}{
		"access_token":  details.AccessToken,
		"refresh_token": details.RefreshToken,
	}, nil
}

func (ai *authInteractor) ValidateAccessToken(ctx context.Context, bearerToken string) (*model.AccessClaims, error) {

	claims, err := ai.jwtConfigurator.GetAccessTokenClaims(bearerToken)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if claims.Authorized == false {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	if claims.Type != model.AccessTokenTypeAuth {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	user, err := ai.UserRepository.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !user.IsActive {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	err = ai.AuthRepository.ValidateAccessToken(ctx, claims.AtID, claims.SessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	return claims, nil
}

func (ai *authInteractor) ValidateTwoFactorAuthToken(ctx context.Context, bearerToken string) (*model.AccessClaims, error) {

	claims, err := ai.jwtConfigurator.GetAccessTokenClaims(bearerToken)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if claims.Authorized == true {
		return nil, fiber.NewError(fiber.StatusBadRequest, "you are already authenticated")
	}

	if claims.Type != model.AccessTokenTypeTwoFactorAuth {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	user, err := ai.UserRepository.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !user.IsActive {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	err = ai.AuthRepository.ValidateAccessToken(ctx, claims.AtID, claims.SessionID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	return claims, nil
}
