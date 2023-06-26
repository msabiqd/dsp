package usecase

import (
	"github.com/SawitProRecruitment/UserService/repository"
)

type UseCase struct {
	Repository repository.RepositoryInterface
}

func NewUseCase(repository repository.RepositoryInterface) *UseCase {
	return &UseCase{Repository: repository}
}
