package controller

type APIController struct {
	Auth            interface{ AuthController }
	User            interface{ UserController }
}
