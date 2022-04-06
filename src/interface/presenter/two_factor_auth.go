package presenter

type twoFactorAuthPresenter struct {
}

type TwoFactorAuthPresenter interface {
}

func NewTwoFactorAuthPresenter() TwoFactorAuthPresenter {
	return &twoFactorAuthPresenter{}
}
