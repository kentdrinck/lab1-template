package service

import (
	"rsoi/internal/model"
)

type Repo interface {
	Add(p *model.Person) error
	Update(p *model.Person) error
	UpdateField(personId int, field string, value any) error
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
	if p.Address != "" {
		if err := s.repo.UpdateField(p.ID, "address", p.Address); err != nil {
			return err
		}
	}
	if p.Name != "" {
		if err := s.repo.UpdateField(p.ID, "name", p.Name); err != nil {
			return err
		}
	}
	if p.Work != "" {
		if err := s.repo.UpdateField(p.ID, "work", p.Work); err != nil {
			return err
		}
	}
	if p.Age != 0 {
		if err := s.repo.UpdateField(p.ID, "age", p.Age); err != nil {
			return err
		}
	}
	return nil
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
