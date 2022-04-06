package registry

import (
	interfaceRepository "auth-project/src/interface/repository"
	usecaseRepository "auth-project/src/usecase/repository"
)

func (r *registry) NewSessionRepository() usecaseRepository.SessionRepository {
	return interfaceRepository.NewSessionRepository(r.db)
}
