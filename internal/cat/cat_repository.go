package cat

import "gorm.io/gorm"

type Repository interface {
	Create(cat *Cat) error
	FindByID(id uint) (*Cat, error)
	List() ([]Cat, error)
	Update(cat *Cat) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new cat repository with the given GORM DB instance.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new Cat record into the database.
func (r *repository) Create(cat *Cat) error {
	return r.db.Create(cat).Error
}

// FindByID retrieves a Cat by its primary key (ID).
func (r *repository) FindByID(id uint) (*Cat, error) {
	var c Cat
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// List retrieves all Cat records from the database.
func (r *repository) List() ([]Cat, error) {
	var cats []Cat
	if err := r.db.Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

// Update applies changes to an existing Cat record in the database.
func (r *repository) Update(cat *Cat) error {
	return r.db.Save(cat).Error
}

// Delete removes a Cat record from the database by its ID.
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Cat{}, id).Error
}
