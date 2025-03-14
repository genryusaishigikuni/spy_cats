package note

import (
	"errors"
	"github.com/genryusaishigikuni/spy_cats/internal/target"
)

type Service interface {
	CreateNote(targetID uint, content string) (*Note, error)
	UpdateNote(noteID uint, content string) (*Note, error)
}

type service struct {
	noteRepo    Repository
	targetRepo  target.Repository
	missionRepo target.Repository
}

func NewService(nRepo Repository, tRepo target.Repository) Service {
	return &service{
		noteRepo:   nRepo,
		targetRepo: tRepo,
	}
}

// CreateNote creates a new note for a target, disallowing creation if the target
// is completed (frozen).
func (s *service) CreateNote(targetID uint, content string) (*Note, error) {
	t, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	if t.Status == "COMPLETED" {
		return nil, errors.New("cannot add note to a completed target")
	}

	// Also check mission status if you want to block notes if mission is completed
	m, err := s.missionRepo.FindByID(t.MissionID)
	if err != nil {
		return nil, errors.New("mission not found")
	}
	if m.Status == "COMPLETED" {
		return nil, errors.New("cannot add note to a completed mission")
	}

	n := &Note{
		TargetID: targetID,
		Content:  content,
	}

	if err := s.noteRepo.Create(n); err != nil {
		return nil, err
	}
	return n, nil
}

// UpdateNote updates an existing note's content, disallowing changes if
// its target or mission is completed.
func (s *service) UpdateNote(noteID uint, content string) (*Note, error) {
	// Find existing note
	n, err := s.noteRepo.FindByID(noteID)
	if err != nil {
		return nil, err
	}

	// Check if note's target is completed
	t, err := s.targetRepo.FindByID(n.TargetID)
	if err != nil {
		return nil, err
	}
	if t.Status == "COMPLETED" {
		return nil, errors.New("cannot update note for a completed target")
	}

	// Check if mission is completed
	m, err := s.missionRepo.FindByID(t.MissionID)
	if err != nil {
		return nil, errors.New("mission not found")
	}
	if m.Status == "COMPLETED" {
		return nil, errors.New("cannot update note for a completed mission")
	}

	// Update note content
	n.Content = content
	if err := s.noteRepo.Update(n); err != nil {
		return nil, err
	}
	return n, nil
}
