package repo

import (
	"errors"
	"rsoi/internal/model"
)

var (
	ErrNotFound = errors.New("not found")
	ErrNilPtr   = errors.New("nil ptr")
)

type InMemoryRepo struct {
	persons []model.Person
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		persons: make([]model.Person, 0),
	}
}

func (r *InMemoryRepo) GetAll() ([]model.Person, error) {
	return r.persons[:len(r.persons)], nil
}

func (r *InMemoryRepo) GetById(id int) (*model.Person, error) {
	for _, p := range r.persons {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, ErrNotFound
}

func (r *InMemoryRepo) Add(p *model.Person) error {
	if p == nil {
		return errors.New("nil ptr")
	}
	p.ID = len(r.persons)
	r.persons = append(r.persons, *p)
	return nil
}

func (r *InMemoryRepo) Update(p *model.Person) error {
	if p == nil {
		return ErrNilPtr
	}
	pp, err := r.GetById(p.ID)
	if err != nil {
		return err
	}
	*pp = *p
	return nil
}

func (r *InMemoryRepo) Delete(id int) error {
	removeIndex := -1
	for i, p := range r.persons {
		if p.ID == id {
			removeIndex = i
		}
	}
	if removeIndex == -1 {
		return ErrNotFound
	}
	r.persons[removeIndex] = r.persons[len(r.persons)-1]
	r.persons = r.persons[:len(r.persons)-1]
	return nil
}
