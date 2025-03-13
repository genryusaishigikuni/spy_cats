package mission

import "gorm.io/gorm"

type Repository interface {
	Create(m *Mission) error
	FindByID(id uint) (*Mission, error)
	Update(m *Mission) error
	Delete(id uint) error
	List() ([]Mission, error)
	FindOngoingByCatID(catID uint) (*Mission, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new Mission record.
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

// Update applies changes to an existing Mission record.
func (r *repository) Update(m *Mission) error {
	return r.db.Save(m).Error
}

// Delete removes a Mission by its ID.
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Mission{}, id).Error
}

// List returns all missions.
func (r *repository) List() ([]Mission, error) {
	var missions []Mission
	if err := r.db.Find(&missions).Error; err != nil {
		return nil, err
	}
	return missions, nil
}

// FindOngoingByCatID checks if the cat has an ongoing mission (Status != 'COMPLETED').
func (r *repository) FindOngoingByCatID(catID uint) (*Mission, error) {
	var m Mission
	if err := r.db.
		Where("cat_id = ? AND status <> ?", catID, "COMPLETED").
		First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
