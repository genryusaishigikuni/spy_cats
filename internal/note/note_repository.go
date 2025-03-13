package note

import "gorm.io/gorm"

type Repository interface {
	Create(n *Note) error
	FindByID(id uint) (*Note, error)
	Update(n *Note) error
	// etc...
}

// repository implements the Repository interface for notes.
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new note repository with the given GORM DB instance.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new Note record into the database.
func (r *repository) Create(n *Note) error {
	return r.db.Create(n).Error
}

// FindByID retrieves a Note by its primary key (ID).
func (r *repository) FindByID(id uint) (*Note, error) {
	var note Note
	if err := r.db.First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

// Update applies changes to an existing Note record in the database.
func (r *repository) Update(n *Note) error {
	return r.db.Save(n).Error
}
