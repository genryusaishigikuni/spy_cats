package mission

import "gorm.io/gorm"

type Repository interface {
	Create(m *Mission) error
	FindByID(id uint) (*Mission, error)
	Update(m *Mission) error
	// add more as needed (List, Delete, etc.)
}
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new mission repository with the given GORM DB instance.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new Mission record into the database.
func (r *repository) Create(m *Mission) error {
	return r.db.Create(m).Error
}

// FindByID retrieves a Mission by its primary key (ID).
func (r *repository) FindByID(id uint) (*Mission, error) {
	var mission Mission
	if err := r.db.First(&mission, id).Error; err != nil {
		return nil, err
	}
	return &mission, nil
}

// Update applies changes to an existing Mission record in the database.
func (r *repository) Update(m *Mission) error {
	return r.db.Save(m).Error
}
