package repositories

import (
	"context"

	"github.com/alex-bezverkhniy/bp-engine/internal/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	ProcessRepository interface {
		Create(ctx context.Context, process *model.Process) (string, error)
		GetByUUID(ctx context.Context, code string, uuid string) (*model.Process, error)
		GetByCode(ctx context.Context, code string, page int, pageSize int) ([]model.Process, error)
		SetStatus(ctx context.Context, code string, uuid string, status string, metadata datatypes.JSON) error
	}
	ProcessRepo struct {
		db *gorm.DB
	}
)

func NewProcessRepository(db *gorm.DB) ProcessRepository {
	return &ProcessRepo{
		db: db,
	}
}

func (r *ProcessRepo) Create(ctx context.Context, process *model.Process) (string, error) {

	if len(process.UUID) == 0 {
		process.UUID = uuid.NewString()
	}
	err := r.db.WithContext(ctx).Create(process).Error
	if err != nil {
		return "", err
	}

	return process.UUID, nil
}

func (r *ProcessRepo) GetByUUID(ctx context.Context, code string, uuid string) (*model.Process, error) {
	var process model.Process
	err := r.db.WithContext(ctx).
		Model(&model.Process{}).
		Preload("CurrentStatus", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Statuses", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Where("code = ? AND uuid = ?", code, uuid).
		First(&process).Error

	return &process, err
}

func (r *ProcessRepo) GetByCode(ctx context.Context, code string, page int, pageSize int) ([]model.Process, error) {
	offset := (page - 1) * pageSize

	var processes []model.Process
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Model(&model.Process{}).
		Preload("Statuses").
		Find(&processes, "code = ?", code).Error
	if err != nil {
		return nil, err
	}

	if len(processes) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return processes, nil

}

func (r *ProcessRepo) SetStatus(ctx context.Context, code string, uuid string, status string, metadata datatypes.JSON) error {
	process, err := r.GetByUUID(ctx, code, uuid)
	if err != nil {
		return err
	}

	newStatus := &model.ProcessStatus{
		ProcessID: process.ID,
		Name:      status,
		Payload:   metadata,
	}

	err = r.db.WithContext(ctx).Model(&model.ProcessStatus{}).Save(&newStatus).Error
	if err != nil {
		return err
	}

	// process.Statuses = append(process.Statuses, *newStatus)
	return r.db.WithContext(ctx).
		Model(&model.Process{}).
		Where("id = ?", process.ID).
		Association("Statuses").
		Append(&newStatus)
}
