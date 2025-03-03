package services

import (
	"user-service/repositories"
	services "user-service/services/user"
)

type Registry struct {
	repository repositories.IRepositoryRegistry
}

type IServiceRegistry interface {
	GetUser() services.IUserService
}

func NewServiceRegistry(repoitory repositories.IRepositoryRegistry) IServiceRegistry {
	return &Registry{repository: repoitory}
}

func (r *Registry) GetUser() services.IUserService {
	return services.NewUserService(r.repository)
}
