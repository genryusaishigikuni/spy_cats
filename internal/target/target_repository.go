package target

import "gorm.io/gorm"

type Repository interface {
	Create(t *Target) error
	FindByID(id uint) (*Target, error)
	FindByMissionID(missionID uint) ([]Target, error)
	Update(t *Target) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(t *Target) error {
	return r.db.Create(t).Error
}

func (r *repository) FindByID(id uint) (*Target, error) {
	var tgt Target
	if err := r.db.First(&tgt, id).Error; err != nil {
		return nil, err
	}
	return &tgt, nil
}

func (r *repository) FindByMissionID(missionID uint) ([]Target, error) {
	var targets []Target
	if err := r.db.Where("mission_id = ?", missionID).Find(&targets).Error; err != nil {
		return nil, err
	}
	return targets, nil
}

func (r *repository) Update(t *Target) error {
	return r.db.Save(t).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Target{}, id).Error
}
