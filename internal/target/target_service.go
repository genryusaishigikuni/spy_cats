package target

import (
	"errors"

	"gorm.io/gorm"
)

// Service defines business operations for the target domain.
type Service interface {
	AddTarget(missionID uint, name string) error
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
func (s *service) AddTarget(missionID uint, name string) error {
	// Retrieve existing targets for the mission
	targets, err := s.repo.FindByMissionID(missionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Enforce max 3 targets
	if len(targets) >= 3 {
		return errors.New("cannot add more than 3 targets to a mission")
	}

	// Create the new target
	t := &Target{
		MissionID: missionID,
		Name:      name,
		Status:    "ONGOING",
	}

	return s.repo.Create(t)
}

// RemoveTarget removes a target by its ID.
func (s *service) RemoveTarget(id uint) error {
	// 1) Check existence
	t, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 2) (Optional) Validate if safe to remove (e.g., mission completed or not)
	// Currently no validation, but I will add it in future

	// 3) Remove the target using the repository's Delete method
	if err := s.repo.Delete(t.ID); err != nil {
		return err
	}

	return nil
}
