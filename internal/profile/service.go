package profile

import (
	"database/sql"
	"errors"
	"log"
)

type Service interface {
	GetProfileService(id int) (*Profile, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

var ErrNotFound = errors.New("not found")

func (s *service) GetProfileService(id int) (*Profile, error) {
	profil, err := s.repo.GetProfileRepository(id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Printf("Error get profile from repository. %v", err)
		return nil, err
	}

	return profil, nil
}
