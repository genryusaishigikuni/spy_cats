package mission

import (
	"errors"
	"time"

	"github.com/genryusaishigikuni/spy_cats/internal/cat"
	"github.com/genryusaishigikuni/spy_cats/internal/target"
)

// Service defines business operations for the mission domain.
type Service interface {
	CreateMission(catID uint, targetNames []string) (*Mission, error)
	CompleteTarget(targetID uint) error
	// Add other methods as needed (ListMissions, etc.).
}

type service struct {
	missionRepo Repository
	catRepo     cat.Repository
	targetRepo  target.Repository
}

// NewService constructs a new mission service with the required repositories.
func NewService(
	mRepo Repository,
	cRepo cat.Repository,
	tRepo target.Repository,
) Service {
	return &service{
		missionRepo: mRepo,
		catRepo:     cRepo,
		targetRepo:  tRepo,
	}
}

// CreateMission creates a new mission, ensuring valid cat and 1–3 target names.
func (s *service) CreateMission(catID uint, targetNames []string) (*Mission, error) {
	// Validate the cat
	_, err := s.catRepo.FindByID(catID)
	if err != nil {
		return nil, errors.New("cat not found")
	}

	// Validate targets (1 to 3)
	if len(targetNames) < 1 || len(targetNames) > 3 {
		return nil, errors.New("mission must have between 1 and 3 targets")
	}

	// Create mission
	m := &Mission{
		CatID:  catID,
		Status: "ONGOING",
	}
	if err := s.missionRepo.Create(m); err != nil {
		return nil, err
	}

	// Create targets for this mission
	for _, tName := range targetNames {
		t := &target.Target{
			MissionID: m.ID,
			Name:      tName,
			Status:    "ONGOING",
		}
		if err := s.targetRepo.Create(t); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// CompleteTarget marks a target as completed and, if all targets in the mission
// are completed, marks the mission as completed as well.
func (s *service) CompleteTarget(targetID uint) error {
	// Fetch the target
	t, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		return err
	}

	// If already completed, return early (optional)
	if t.Status == "COMPLETED" {
		return errors.New("target is already completed")
	}

	// Mark the target as completed
	t.Status = "COMPLETED"
	now := time.Now()
	t.CompletedAt = &now

	if err := s.targetRepo.Update(t); err != nil {
		return err
	}

	// Check whether all targets for this mission are completed
	targets, err := s.targetRepo.FindByMissionID(t.MissionID)
	if err != nil {
		return err
	}

	allDone := true
	for _, each := range targets {
		if each.Status != "COMPLETED" {
			allDone = false
			break
		}
	}

	// If all done, update the mission as completed
	if allDone {
		mission, err := s.missionRepo.FindByID(t.MissionID)
		if err != nil {
			return err
		}

		mission.Status = "COMPLETED"
		mission.CompletedAt = &now

		if err := s.missionRepo.Update(mission); err != nil {
			return err
		}
	}

	return nil
}
