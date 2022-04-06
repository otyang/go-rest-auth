package registry

import (
	interfaceController "auth-project/src/interface/controller"
	interfacePresenter "auth-project/src/interface/presenter"
	interfaceRepository "auth-project/src/interface/repository"
	usecaseInteractor "auth-project/src/usecase/interactor"
	usecasePresenter "auth-project/src/usecase/presenter"
	usecaseRepository "auth-project/src/usecase/repository"
)

func (r *registry) NewAuthController() interfaceController.AuthController {
	return interfaceController.NewAuthController(r.NewAuthInteractor())
}

func (r *registry) NewAuthInteractor() usecaseInteractor.AuthInteractor {
	return usecaseInteractor.NewAuthInteractor(r.NewAuthRepository(), r.NewSessionRepository(), r.NewUserRepository(), r.NewTokenRepository(), r.NewAuthPresenter(), r.jwtConf)
}

func (r *registry) NewAuthRepository() usecaseRepository.AuthRepository {
	return interfaceRepository.NewAuthRepository(r.db, r.rdb)
}

func (r *registry) NewAuthPresenter() usecasePresenter.AuthPresenter {
	return interfacePresenter.NewAuthPresenter()
}
