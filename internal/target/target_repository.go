package target

import "gorm.io/gorm"

// Repository defines the methods to interact with the Targets in the database.
type Repository interface {
	Create(t *Target) error
	FindByID(id uint) (*Target, error)
	FindByMissionID(missionID uint) ([]Target, error)
	Update(t *Target) error
	Delete(id uint) error
}

// repository is the concrete implementation of the Repository interface.
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new target repository with the given GORM DB instance.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new Target record into the database.
func (r *repository) Create(t *Target) error {
	return r.db.Create(t).Error
}

// FindByID retrieves a Target by its primary key (ID).
func (r *repository) FindByID(id uint) (*Target, error) {
	var tgt Target
	if err := r.db.First(&tgt, id).Error; err != nil {
		return nil, err
	}
	return &tgt, nil
}

// FindByMissionID retrieves all Targets associated with a particular mission.
func (r *repository) FindByMissionID(missionID uint) ([]Target, error) {
	var targets []Target
	if err := r.db.Where("mission_id = ?", missionID).Find(&targets).Error; err != nil {
		return nil, err
	}
	return targets, nil
}

// Update applies changes to an existing Target record in the database.
func (r *repository) Update(t *Target) error {
	return r.db.Save(t).Error
}

// Delete removes a Target record from the database by its ID.
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Target{}, id).Error
}
