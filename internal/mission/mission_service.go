package mission

import (
	"errors"
	"fmt"
	"time"

	"github.com/genryusaishigikuni/spy_cats/internal/cat"
	"github.com/genryusaishigikuni/spy_cats/internal/target"
	"gorm.io/gorm"
)

type Service interface {
	CreateMission(catID uint, targetNames []string) (*Mission, error)
	CompleteTarget(targetID uint) error

	ListMissions() ([]Mission, error)
	GetMissionByID(id uint) (*Mission, error)
	DeleteMission(id uint) error
	MarkMissionComplete(missionID uint) error
	AssignCat(missionID, catID uint) error

	// AddTargetToMission New: Add a target to an existing mission.
	AddTargetToMission(missionID uint, name, country, notes string) error
}

type service struct {
	missionRepo Repository
	catRepo     cat.Repository
	targetRepo  target.Repository
}

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

// CreateMission creates a new mission, ensuring the cat is valid and doesn't have an ongoing mission, plus 1–3 targets.
func (s *service) CreateMission(catID uint, targetNames []string) (*Mission, error) {
	// Validate the cat
	_, err := s.catRepo.FindByID(catID)
	if err != nil {
		return nil, errors.New("cat not found")
	}

	// Check if the cat already has an ongoing mission
	_, err = s.missionRepo.FindOngoingByCatID(catID)
	if err == nil {
		return nil, errors.New("this cat already has an ongoing mission")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
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

	// Create targets for this mission using the target repository directly.
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

// AddTargetToMission adds a new target to an existing mission,
// ensuring the mission is ongoing and does not exceed 3 targets.
func (s *service) AddTargetToMission(missionID uint, name, country, notes string) error {
	// 1) Check mission exists and is not completed
	m, err := s.missionRepo.FindByID(missionID)
	if err != nil {
		return fmt.Errorf("mission not found: %w", err)
	}
	if m.Status == "COMPLETED" {
		return errors.New("cannot add target to a completed mission")
	}

	// 2) Check number of existing targets
	existingTargets, err := s.targetRepo.FindByMissionID(missionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if len(existingTargets) >= 3 {
		return errors.New("cannot add more than 3 targets to a mission")
	}

	// 3) Create the new target with additional fields Country and Notes
	t := &target.Target{
		MissionID: missionID,
		Name:      name,
		Country:   country,
		Notes:     notes,
		Status:    "ONGOING",
	}
	return s.targetRepo.Create(t)
}

// CompleteTarget marks a target as completed and, if all targets in the mission are completed,
// marks the mission as completed as well.
func (s *service) CompleteTarget(targetID uint) error {
	// Fetch the target
	t, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		return err
	}

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

	// Check if all targets for this mission are completed
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

	if allDone {
		if err := s.markMissionCompleted(t.MissionID); err != nil {
			return err
		}
	}

	return nil
}

// ListMissions returns all missions.
func (s *service) ListMissions() ([]Mission, error) {
	return s.missionRepo.List()
}

// GetMissionByID returns a single mission by ID.
func (s *service) GetMissionByID(id uint) (*Mission, error) {
	m, err := s.missionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("mission not found")
	}
	return m, nil
}

// DeleteMission removes a mission if it isn't assigned to a cat.
func (s *service) DeleteMission(id uint) error {
	m, err := s.missionRepo.FindByID(id)
	if err != nil {
		return errors.New("mission not found")
	}

	// If a cat is assigned, forbid deletion.
	if m.CatID != 0 {
		return errors.New("cannot delete a mission that is assigned to a cat")
	}

	return s.missionRepo.Delete(id)
}

// MarkMissionComplete forcibly completes a mission.
func (s *service) MarkMissionComplete(missionID uint) error {
	return s.markMissionCompleted(missionID)
}

// AssignCat assigns a cat to an existing mission if valid.
func (s *service) AssignCat(missionID, catID uint) error {
	m, err := s.missionRepo.FindByID(missionID)
	if err != nil {
		return errors.New("mission not found")
	}
	if m.Status == "COMPLETED" {
		return errors.New("cannot assign a cat to a completed mission")
	}

	_, err = s.catRepo.FindByID(catID)
	if err != nil {
		return errors.New("cat not found")
	}

	// Check if the cat is free.
	_, err = s.missionRepo.FindOngoingByCatID(catID)
	if err == nil {
		return errors.New("this cat is already on another ongoing mission")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	m.CatID = catID
	return s.missionRepo.Update(m)
}

// markMissionCompleted is an internal helper to mark a mission as completed.
func (s *service) markMissionCompleted(missionID uint) error {
	m, err := s.missionRepo.FindByID(missionID)
	if err != nil {
		return errors.New("mission not found")
	}

	now := time.Now()
	m.Status = "COMPLETED"
	m.CompletedAt = &now

	return s.missionRepo.Update(m)
}
