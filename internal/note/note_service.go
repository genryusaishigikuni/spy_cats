package note

import (
	"errors"

	"github.com/genryusaishigikuni/spy_cats/internal/target"
)

// Service defines business operations for the note domain.
type Service interface {
	CreateNote(targetID uint, content string) (*Note, error)
	UpdateNote(noteID uint, content string) (*Note, error)
}

type service struct {
	noteRepo   Repository
	targetRepo target.Repository
}

// NewService constructs a new note service with the required repositories.
func NewService(nRepo Repository, tRepo target.Repository) Service {
	return &service{
		noteRepo:   nRepo,
		targetRepo: tRepo,
	}
}

// CreateNote creates a new note for a target, disallowing creation if the target
// is completed (frozen).
func (s *service) CreateNote(targetID uint, content string) (*Note, error) {
	// Check if target is completed
	t, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		return nil, err
	}

	if t.Status == "COMPLETED" {
		return nil, errors.New("cannot add note to a completed target")
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
// its target is completed.
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

	// Update note content
	n.Content = content
	if err := s.noteRepo.Update(n); err != nil {
		return nil, err
	}
	return n, nil
}
