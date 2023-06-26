package handler

import (
	"github.com/SawitProRecruitment/UserService/usecase"
)

type Server struct {
	UseCase usecase.UseCaseInterface
}

// type NewServerOptions struct {
// 	Repository repository.RepositoryInterface
// }

func NewServer(usecase usecase.UseCaseInterface) *Server {
	return &Server{UseCase: usecase}
}
