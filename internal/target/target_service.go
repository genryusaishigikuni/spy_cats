package target

import (
	"fmt"
)

// Service defines business operations for the target domain.
type Service interface {
	RemoveTarget(id uint) error
}

type service struct {
	repo Repository
}

// NewService constructs a new target service with the required repository.
func NewService(r Repository) Service {
	return &service{repo: r}
}

// AddTarget adds a new target to the given mission, ensuring the mission
// does not exceed 3 targets.

// RemoveTarget removes a target by its ID.
// RemoveTarget removes a target by its ID.
func (s *service) RemoveTarget(id uint) error {
	// 1) Check existence
	t, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 2) Validate if the target is completed
	if t.Status == "COMPLETED" {
		return fmt.Errorf("cannot delete a completed target")
	}

	// 3) Remove the target using the repository's Delete method
	if err := s.repo.Delete(t.ID); err != nil {
		return err
	}

	return nil
}
