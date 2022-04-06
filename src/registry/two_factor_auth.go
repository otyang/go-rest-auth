package registry

import (
	interfaceController "auth-project/src/interface/controller"
	interfacePresenter "auth-project/src/interface/presenter"
	interfaceRepository "auth-project/src/interface/repository"
	usecaseInteractor "auth-project/src/usecase/interactor"
	usecasePresenter "auth-project/src/usecase/presenter"
	usecaseRepository "auth-project/src/usecase/repository"
)

func (r *registry) NewTwoFactorAuthController() interfaceController.TwoFactorAuthController {
	return interfaceController.NewTwoFactorAuthController(r.NewTwoFactorAuthInteractor())
}

func (r *registry) NewTwoFactorAuthInteractor() usecaseInteractor.TwoFactorAuthInteractor {
	return usecaseInteractor.NewTwoFactorAuthInteractor(r.NewAuthRepository(), r.NewSessionRepository(),
		r.NewTwoFactorAuthRepository(), r.NewUserRepository(), r.NewTokenRepository(), r.NewTwoFactorAuthPresenter(),
		r.jwtConf)
}

func (r *registry) NewTwoFactorAuthRepository() usecaseRepository.TwoFactorAuthRepository {
	return interfaceRepository.NewTwoFactorAuthRepository(r.db)
}

func (r *registry) NewTwoFactorAuthPresenter() usecasePresenter.TwoFactorAuthPresenter {
	return interfacePresenter.NewTwoFactorAuthPresenter()
}
