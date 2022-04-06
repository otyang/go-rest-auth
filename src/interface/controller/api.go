package controller

type APIController struct {
	Auth          interface{ AuthController }
	QrCodeAuth    interface{ QrCodeAuthController }
	TwoFactorAuth interface{ TwoFactorAuthController }
	Token         interface{ TokenController }
	User          interface{ UserController }
}
