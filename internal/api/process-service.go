package api

import (
	"context"

	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	ProcessService interface {
		Submit(ctx context.Context, process *ProcessDTO) (string, error)
		Get(ctx context.Context, code string, uuid string, page int, pageSize int) (ProcessListDTO, error)
		AssignStatus(ctx context.Context, code string, uuid string, status string, metadata Payload) error
	}
	ProcessSrvc struct {
		repo ProcessRepository
	}
)

var (
	ErrProcessNotFound     error = errors.New("process not found")
	ErrCannotCreateProcess error = errors.New("cannot create process")
)

func NewProcessService(repo ProcessRepository) ProcessService {
	return &ProcessSrvc{
		repo: repo,
	}
}

func (s *ProcessSrvc) Submit(ctx context.Context, process *ProcessDTO) (string, error) {
	uuid, err := s.repo.Create(ctx, process.toEntity())
	if err != nil {
		return "", errors.Join(err, ErrCannotCreateProcess)
	}
	return uuid, nil
}

func (s *ProcessSrvc) Get(ctx context.Context, code string, uuid string, page int, pageSize int) (ProcessListDTO, error) {
	var processes ProcessList
	var process *Process
	var err error

	// Get by code
	if len(uuid) == 0 {
		if page <= 0 {
			page = DEFAULT_PAGE
		}

		if pageSize <= 0 {
			pageSize = DEFAULT_PAGE_SIZE
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
			return nil, ErrProcessNotFound
		}

		return nil, err
	}

	return processes.toDTO(), nil
}

func (s *ProcessSrvc) AssignStatus(ctx context.Context, code string, uuid string, status string, metadata Payload) error {
	err := s.repo.SetStatus(ctx, code, uuid, status, datatypes.JSON(metadata.toBytes()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProcessNotFound
		}
		return err
	}

	return nil
}
