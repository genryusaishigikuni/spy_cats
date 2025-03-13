package cat

import "time"

type Cat struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Breed     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
