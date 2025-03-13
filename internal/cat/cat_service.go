package cat

import (
	"errors"
)

// Service defines business operations for the cat domain.
type Service interface {
	CreateCat(name, breed string) (*Cat, error)
	GetCat(id uint) (*Cat, error)
	ListCats() ([]Cat, error)
	UpdateCat(id uint, name, breed string) (*Cat, error)
	DeleteCat(id uint) error
}

type service struct {
	repo Repository
}

// NewService creates a new cat service with the given cat repository.
func NewService(r Repository) Service {
	return &service{repo: r}
}

// CreateCat creates a new Cat record after basic validations.
func (s *service) CreateCat(name, breed string) (*Cat, error) {
	if name == "" || breed == "" {
		return nil, errors.New("cat name and breed must not be empty")
	}
	c := &Cat{Name: name, Breed: breed}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

// GetCat retrieves a cat by its ID.
func (s *service) GetCat(id uint) (*Cat, error) {
	// Call the repository's FindByID method
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("cat not found") // Return custom error if not found
	}
	return cat, nil
}

// ListCats retrieves all cats.
func (s *service) ListCats() ([]Cat, error) {
	// Retrieve all cats from the repository
	cats, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return cats, nil
}

// UpdateCat updates an existing cat's data.
func (s *service) UpdateCat(id uint, name, breed string) (*Cat, error) {
	// Find the cat by ID first
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("cat not found") // Returns custom error if cat not found
	}

	// Updates the cat's details
	c.Name = name
	c.Breed = breed

	// Calls the repository to update the cat
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}

	return c, nil
}

// DeleteCat removes a cat from the database by ID.
func (s *service) DeleteCat(id uint) error {
	// Call the repository's Delete method
	if err := s.repo.Delete(id); err != nil {
		return errors.New("could not delete cat") // Returns custom error if delete fails
	}
	return nil
}
