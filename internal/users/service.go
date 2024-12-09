package users

import (
	"database/sql"
	"errors"
	"log"

	"github.com/DanyAdhi/learn-golang/internal/utils"
)

type Service interface {
	GetAllUsersService(params GetAllUsersParmas) (*GetAllUsersResponse, error)
	GetOneUsersService(id int) (*User, error)
	CreateUsersService(user *Createuser) error
	UpdateUsersService(id int, user *UpdateUser) error
	DeleteUsersService(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrGeneratePassword = errors.New("failed generate password")
var ErrUserNotFound = errors.New("user not found")

func (s *service) GetAllUsersService(p GetAllUsersParmas) (*GetAllUsersResponse, error) {
	users, totalRecords, err := s.repo.GetAllUsersRepository(p)
	if err != nil {
		log.Printf("Error get user from repository. %v", err)
		return nil, err
	}

	pagination := utils.GeneratePagination(p.Limit, p.Page, totalRecords)

	response := &GetAllUsersResponse{
		Users:      users,
		Pagination: pagination,
	}

	return response, nil
}

func (s *service) GetOneUsersService(id int) (*User, error) {
	user, err := s.repo.GetOneUsersRepository(id)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		log.Printf("Error get user form repository. %v", err)
		return nil, err
	}

	return user, nil
}

func (s *service) CreateUsersService(user *Createuser) error {
	// check email
	exist, err := s.repo.CheckEmailExists(user.Email)
	if err != nil {
		return err
	}
	if exist {
		return ErrEmailAlreadyExists
	}

	password, err := utils.BcryptHashPassword("password")
	if err != nil {
		log.Printf("error generate password. %v", err)
		return ErrGeneratePassword
	}

	err = s.repo.StoreUsersRepository(user, password)
	if err != nil {
		log.Printf("Error create user %v", err)
		return err
	}
	return nil
}

func (s *service) UpdateUsersService(id int, user *UpdateUser) error {
	_, err := s.repo.GetOneUsersRepository(id)
	if err == sql.ErrNoRows {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	err = s.repo.UpdateUsersRepository(id, user)
	if err != nil {
		log.Printf("Error update user. %v", err)
		return err
	}
	return nil
}

func (s *service) DeleteUsersService(id int) error {
	err := s.repo.DeleteUsersRepository(id)
	if err != nil {
		return err
	}
	return nil
}
