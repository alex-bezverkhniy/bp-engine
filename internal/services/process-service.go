package services

import (
	"context"

	"github.com/alex-bezverkhniy/bp-engine/internal/model"
	"github.com/alex-bezverkhniy/bp-engine/internal/repositories"
	"github.com/alex-bezverkhniy/bp-engine/pkg/engine"
	"github.com/alex-bezverkhniy/bp-engine/pkg/engine/config"
	"github.com/alex-bezverkhniy/bp-engine/pkg/engine/handlers"
	"github.com/alex-bezverkhniy/bp-engine/pkg/engine/validators"

	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	ProcessService interface {
		Submit(ctx context.Context, process *model.ProcessDTO) (string, error)
		Get(ctx context.Context, code string, uuid string, page int, pageSize int) (model.ProcessListDTO, error)
		AssignStatus(ctx context.Context, code string, uuid string, status string, metadata model.Payload) error
	}
	ProcessSrvc struct {
		validator      validators.Validator
		statusHandlers handlers.StatusHandlers
		repo           repositories.ProcessRepository
	}
)

func NewProcessService(repo repositories.ProcessRepository, validator validators.Validator, statusHandlers handlers.StatusHandlers) ProcessService {
	return &ProcessSrvc{
		validator:      validator,
		statusHandlers: statusHandlers,
		repo:           repo,
	}
}

func (s *ProcessSrvc) Submit(ctx context.Context, process *model.ProcessDTO) (string, error) {
	uuid, err := s.repo.Create(ctx, process.ToEntity())
	if err != nil {
		return "", errors.Join(err, engine.ErrCannotCreateProcess)
	}
	return uuid, nil
}

func (s *ProcessSrvc) Get(ctx context.Context, code string, uuid string, page int, pageSize int) (model.ProcessListDTO, error) {
	var processes model.ProcessList
	var process *model.Process
	var err error

	// Get by code
	if len(uuid) == 0 {
		if page <= 0 {
			page = config.DEFAULT_PAGE
		}

		if pageSize <= 0 {
			pageSize = config.DEFAULT_PAGE_SIZE
		}
		processes, err = s.repo.GetByCode(ctx, code, page, pageSize)
	} else {
		process, err = s.repo.GetByUUID(ctx, code, uuid)
		if err == nil {
			processes = append(processes, *process)
		}
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, engine.ErrProcessNotFound
		}

		return nil, err
	}

	return processes.ToDTO(), nil
}

func (s *ProcessSrvc) AssignStatus(ctx context.Context, code string, uuid string, status string, payload model.Payload) error {
	// Check process exist
	process, err := s.repo.GetByUUID(ctx, code, uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return engine.ErrProcessNotFound
		}

		return err
	}
	newStatus := model.ProcessStatusDTO{
		Name:    status,
		Payload: payload,
	}
	// Validate the status
	err = s.validator.Validate(process.ToDTO(), newStatus)
	if err != nil {
		return err
	}

	// Handle the status changes

	err = s.repo.SetStatus(ctx, code, uuid, status, datatypes.JSON(payload.ToBytes()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return engine.ErrProcessNotFound
		}
		return err
	}

	return nil
}
