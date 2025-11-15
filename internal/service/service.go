package service

import (
	"rsoi/internal/model"
)

type Repo interface {
	Add(p *model.Person) error
	Update(p *model.Person) error
	Delete(id int) error
	GetAll() ([]model.Person, error)
	GetById(id int) (*model.Person, error)
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(p *model.Person) error {
	return s.repo.Add(p)
}

func (s *Service) Update(p *model.Person) error {
	return s.repo.Update(p)
}

func (s Service) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s Service) GetAll() ([]model.Person, error) {
	return s.repo.GetAll()
}

func (s Service) GetById(id int) (*model.Person, error) {
	return s.repo.GetById(id)
}
