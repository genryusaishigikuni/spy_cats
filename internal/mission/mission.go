package mission

import (
	"time"
)

// Mission model includes references to CatID, plus a CompletedAt if the mission is done.
type Mission struct {
	ID          uint       `gorm:"primaryKey"`
	CatID       uint       // which cat is assigned
	Status      string     // "ONGOING" or "COMPLETED"
	CompletedAt *time.Time // null if not completed
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
