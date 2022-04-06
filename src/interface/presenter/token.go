package presenter

type tokenPresenter struct {
}

type TokenPresenter interface {
}

func NewTokenPresenter() TokenPresenter {
	return &tokenPresenter{}
}
