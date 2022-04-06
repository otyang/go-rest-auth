package registry

import (
	interfaceController "auth-project/src/interface/controller"
	interfacePresenter "auth-project/src/interface/presenter"
	interfaceRepository "auth-project/src/interface/repository"
	usecaseInteractor "auth-project/src/usecase/interactor"
	usecasePresenter "auth-project/src/usecase/presenter"
	usecaseRepository "auth-project/src/usecase/repository"
)

func (r *registry) NewUserController() interfaceController.UserController {
	return interfaceController.NewUserController(r.NewUserInteractor())
}

func (r *registry) NewUserInteractor() usecaseInteractor.UserInteractor {
	return usecaseInteractor.NewUserInteractor(r.NewUserRepository(), r.NewUserPresenter())
}

func (r *registry) NewUserRepository() usecaseRepository.UserRepository {
	return interfaceRepository.NewUserRepository(r.db, r.rdb)
}

func (r *registry) NewUserPresenter() usecasePresenter.UserPresenter {
	return interfacePresenter.NewUserPresenter()
}
