package note

import "time"

type Note struct {
	ID        uint `gorm:"primaryKey"`
	TargetID  uint
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
