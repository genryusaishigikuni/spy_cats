package mission

import "time"

type Mission struct {
	ID          uint `gorm:"primaryKey"`
	CatID       uint
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}
