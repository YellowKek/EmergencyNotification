package service

import (
	"app/internal/model"
	"app/internal/repository"
)

type UserService struct {
	r repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{r}
}

func (s *UserService) CreateUser(user model.User) error {
	return s.r.Create(user)
}

func (s *UserService) GetByEmail(email string) (*model.User, error) {
	return s.r.GetByEmail(email)
}

func (s *UserService) CreateTables() error {
	return s.r.CreateTable()
}

func (s *UserService) GetAll() ([]*model.User, error) {
	return s.r.GetAll()
}
