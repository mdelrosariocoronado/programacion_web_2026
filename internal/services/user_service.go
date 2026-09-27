package services

import (
	"emprendimientos.com/servidor-go/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

// Nota la N mayúscula
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}