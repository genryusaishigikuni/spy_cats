package target

import "time"

// Target now includes Country and Notes to match the requirement
type Target struct {
	ID          uint `gorm:"primaryKey"`
	MissionID   uint `gorm:"index"` // belongs to a particular mission
	Name        string
	Country     string
	Notes       string
	Status      string // "ONGOING" or "COMPLETED"
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
