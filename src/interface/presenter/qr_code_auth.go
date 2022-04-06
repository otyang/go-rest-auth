package presenter

type qrCodeAuthPresenter struct {
}

type QrCodeAuthPresenter interface {
}

func NewQrCodeAuthPresenter() QrCodeAuthPresenter {
	return &qrCodeAuthPresenter{}
}
