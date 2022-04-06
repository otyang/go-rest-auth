package controller

import (
	"auth-project/src/domain/model"
	"auth-project/src/usecase/interactor"
	"auth-project/tools"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

type LoginReq struct {
	Status int
	Conn   *websocket.Conn
}

var tokenarr = make(map[string]LoginReq)

type authController struct {
	authInteractor  interactor.AuthInteractor
	jwtConfigurator *tools.JwtConfigurator
}

type AuthController interface {
	Authenticate(ctx *fiber.Ctx) error
	GenerateQrCode(ctx *fiber.Ctx) error
	LoginQrCode(ctx *fiber.Ctx) error
	QrPolling(c *websocket.Conn)
	RefreshToken(ctx *fiber.Ctx) error
	ValidateAccessTokenID(ctx *fiber.Ctx) error
}

func NewAuthController(ai interactor.AuthInteractor, jc *tools.JwtConfigurator) AuthController {
	return &authController{ai, jc}
}

// SignIn accepts the user form data and returns a token to authorize a client.
func (ac *authController) Authenticate(ctx *fiber.Ctx) error {

	var err error
	authData := new(model.AuthenticationData)
	err = ctx.BodyParser(&authData)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	usr, err := ac.authInteractor.GetUserByEmailAndPassword(ctx.Context(),
		authData.Email, authData.Password)
	if err != nil {
		err = errors.New("incorrect email or password")
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	usrInfo := &model.UserRedisSessionData{UserID: usr.ID.String()}

	details, err := ac.jwtConfigurator.GenerateTokenPair(usrInfo.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = ac.authInteractor.StoreAuth(ctx.Context(), usrInfo, details)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"access_token":  details.AccessToken,
		"refresh_token": details.RefreshToken,
	})
}

func (ac *authController) GenerateQrCode(ctx *fiber.Ctx) error {

	qrCode, tokenID, err := ac.authInteractor.GenerateQrCode(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"qrCode": qrCode,
		"tokenID": tokenID,
	})
}

func (ac *authController) LoginQrCode(ctx *fiber.Ctx) error {
	token := ctx.Params("token")
	fmt.Println(token)
	fmt.Println(tokenarr[token])
	fmt.Println(tokenarr[token].Conn)
	err := tokenarr[token].Conn.WriteMessage(1, []byte("123"))
	if err != nil {
		return err
	}
	return ctx.SendString(ctx.Params("token"))
}

func (ac *authController) QrPolling(c *websocket.Conn) {
	token := c.Params("token")
	fmt.Println(token)

	tokenarr[token] = LoginReq{Conn: c}
	// c.Locals is added to the *websocket.Conn
	log.Println(c.Locals("allowed"))  // true
	log.Println(c.Params("id"))       // 123
	log.Println(c.Query("v"))         // 1.0
	log.Println(c.Cookies("session")) // ""

	var (
		mt  int
		msg []byte
		err error
	)

	fmt.Println("13213")
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (ac *authController) RefreshToken(ctx *fiber.Ctx) error {

	var err error
	bearerToken := ctx.Get("Authorization")
	if bearerToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Authorization token missing")
	}

	claims, err := ac.jwtConfigurator.GetRefreshTokenClaims(bearerToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if claims.RtID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "refresh token missing")
	}

	UserID, err := ac.authInteractor.FetchAuth(ctx.Context(), claims.RtID)
	if err != nil {
		if err.Error() == "at token dont found" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if UserID != claims.UserID {
		err = errors.New("invalid token")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// All OK, re-generate the new pair and send to client,
	// we could only generate an access token as well.
	details, err := ac.jwtConfigurator.GenerateTokenPair(claims.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	usrInfo := &model.UserRedisSessionData{UserID: claims.UserID}

	err = ac.authInteractor.StoreAuth(ctx.Context(), usrInfo, details)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send the generated token pair to the client.
	// The tokenPair looks like: {"access_token": $token, "refresh_token": $token}
	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"access_token":  details.AccessToken,
		"refresh_token": details.RefreshToken,
	})
}

func (ac *authController) ValidateAccessTokenID(ctx *fiber.Ctx) error {

	bearerToken := ctx.Get("Authorization")
	if bearerToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Authorization token missing")
	}

	claims, err := ac.jwtConfigurator.GetAccessTokenClaims(bearerToken)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// set user id from token in context
	if claims.UserID != "" && claims.AtID != "" {
		ctx.Context().SetUserValue("token_user_id", claims.UserID)
		ctx.Context().SetUserValue("token_at_id", claims.AtID)
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "claims is missing")
	}

	err = ac.authInteractor.ValidateAccessTokenID(ctx.Context(), claims.UserID, claims.AtID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	return nil
}
