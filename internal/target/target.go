package target

import "time"

type Target struct {
	ID          uint `gorm:"primaryKey"`
	MissionID   uint
	Name        string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}
