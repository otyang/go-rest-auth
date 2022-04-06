package presenter

type authPresenter struct {
}

type AuthPresenter interface {
}

func NewAuthPresenter() AuthPresenter {
	return &authPresenter{}
}
