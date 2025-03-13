package cat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Service defines business operations for the cat domain.
type Service interface {
	CreateCat(name, breed string, years int, salary float64) (*Cat, error)
	GetCat(id uint) (*Cat, error)
	ListCats() ([]Cat, error)
	UpdateCat(id uint, name, breed string, years int, salary float64) (*Cat, error)
	DeleteCat(id uint) error
}

type service struct {
	repo Repository
}

// NewService creates a new cat service with the given cat repository.
func NewService(r Repository) Service {
	return &service{repo: r}
}

// CreateCat creates a new Cat record after validations (including breed).
func (s *service) CreateCat(name, breed string, years int, salary float64) (*Cat, error) {
	if name == "" {
		return nil, errors.New("cat name cannot be empty")
	}
	if years < 0 {
		return nil, errors.New("years of experience cannot be negative")
	}
	if salary < 0 {
		return nil, errors.New("salary cannot be negative")
	}
	if err := validateBreed(breed); err != nil {
		return nil, err
	}

	c := &Cat{
		Name:              name,
		Breed:             breed,
		YearsOfExperience: years,
		Salary:            salary,
	}

	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

// GetCat retrieves a cat by its ID.
func (s *service) GetCat(id uint) (*Cat, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("cat not found")
	}
	return cat, nil
}

// ListCats retrieves all cats.
func (s *service) ListCats() ([]Cat, error) {
	return s.repo.List()
}

// UpdateCat updates an existing cat's data, including breed validation.
func (s *service) UpdateCat(id uint, name, breed string, years int, salary float64) (*Cat, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("cat not found")
	}

	if name == "" {
		return nil, errors.New("cat name cannot be empty")
	}
	if years < 0 {
		return nil, errors.New("years of experience cannot be negative")
	}
	if salary < 0 {
		return nil, errors.New("salary cannot be negative")
	}
	if err := validateBreed(breed); err != nil {
		return nil, err
	}

	c.Name = name
	c.Breed = breed
	c.YearsOfExperience = years
	c.Salary = salary

	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

// DeleteCat removes a cat from the database by ID.
func (s *service) DeleteCat(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return errors.New("could not delete cat")
	}
	return nil
}

// validateBreed checks if the given breed exists via TheCatAPI.
func validateBreed(breed string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/breeds", nil)
	if err != nil {
		return fmt.Errorf("could not create request to thecatapi: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not reach thecatapi: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("thecatapi responded with status %d", resp.StatusCode)
	}

	var breeds []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return fmt.Errorf("could not decode breed data: %w", err)
	}

	// Check if breed is in the list
	for _, b := range breeds {
		if strings.EqualFold(b.Name, breed) {
			return nil
		}
	}
	return fmt.Errorf("invalid cat breed: %s", breed)
}
