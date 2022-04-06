package registry

import (
	interfaceController "auth-project/src/interface/controller"
	interfacePresenter "auth-project/src/interface/presenter"
	interfaceRepository "auth-project/src/interface/repository"
	usecaseInteractor "auth-project/src/usecase/interactor"
	usecasePresenter "auth-project/src/usecase/presenter"
	usecaseRepository "auth-project/src/usecase/repository"
)

func (r *registry) NewQrCodeAuthController() interfaceController.QrCodeAuthController {
	return interfaceController.NewQrCodeAuthController(r.NewQrCodeAuthInteractor())
}

func (r *registry) NewQrCodeAuthInteractor() usecaseInteractor.QrCodeAuthInteractor {
	return usecaseInteractor.NewQrCodeAuthInteractor(r.NewAuthRepository(), r.NewSessionRepository(), r.NewUserRepository(), r.NewQrCodeAuthRepository(), r.NewQrCodeAuthPresenter(), r.jwtConf)
}

func (r *registry) NewQrCodeAuthRepository() usecaseRepository.QrCodeAuthRepository {
	return interfaceRepository.NewQrCodeAuthRepository(r.db)
}

func (r *registry) NewQrCodeAuthPresenter() usecasePresenter.QrCodeAuthPresenter {
	return interfacePresenter.NewQrCodeAuthPresenter()
}
