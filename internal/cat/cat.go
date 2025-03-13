package cat

import "time"

type Cat struct {
	ID                uint `gorm:"primaryKey"`
	Name              string
	Breed             string
	YearsOfExperience int
	Salary            float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
