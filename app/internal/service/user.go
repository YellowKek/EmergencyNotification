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

func (s *UserService) GetEmergencyGroups(userID int32) (map[string]string, error) {
	return s.r.GetEmergencyGroups(userID)
}

func (s *UserService) GetByID(userID int32) (*model.User, error) {
	return s.r.GetById(userID)
}

func (s *UserService) AddEmergencyGroup(userID int32, groupType, value string) error {
	return s.r.AddEmergencyGroup(userID, groupType, value)
}
