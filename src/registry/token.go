package registry

import (
	interfaceController "auth-project/src/interface/controller"
	interfacePresenter "auth-project/src/interface/presenter"
	interfaceRepository "auth-project/src/interface/repository"
	usecaseInteractor "auth-project/src/usecase/interactor"
	usecasePresenter "auth-project/src/usecase/presenter"
	usecaseRepository "auth-project/src/usecase/repository"
)

func (r *registry) NewTokenController() interfaceController.TokenController {
	return interfaceController.NewTokenController(r.NewTokenInteractor())
}

func (r *registry) NewTokenInteractor() usecaseInteractor.TokenInteractor {
	return usecaseInteractor.NewTokenInteractor(r.NewTokenRepository(), r.NewUserRepository(), r.NewTokenPresenter())
}

func (r *registry) NewTokenRepository() usecaseRepository.TokenRepository {
	return interfaceRepository.NewTokenRepository(r.db)
}

func (r *registry) NewTokenPresenter() usecasePresenter.TokenPresenter {
	return interfacePresenter.NewTokenPresenter()
}
