package main

import (
	"context"

	"errors"

	"gorm.io/gorm"
)

type (
	ProcessService interface {
		Submit(ctx context.Context, process *ProcessDTO) (string, error)
		Get(ctx context.Context, code string, uuid string) ([]ProcessDTO, error)
		AssignStatus(ctx context.Context, code string, uuid string, status string) error
	}
	ProcessSrvc struct {
		repo ProcessRepository
	}
)

var (
	ErrNoProcessesFound    error = errors.New("no processes found")
	ErrProcessNotFound     error = errors.New("process not found")
	ErrCannotCreateProcess error = errors.New("process not found")
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

func (s *ProcessSrvc) Get(ctx context.Context, code string, uuid string) ([]ProcessDTO, error) {
	var processes ProcessList
	var process *Process
	var err error

	if len(uuid) == 0 {
		processes, err = s.repo.GetByCode(ctx, code)
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

func (s *ProcessSrvc) AssignStatus(ctx context.Context, code string, uuid string, status string) error {
	err := s.repo.SetStatus(ctx, code, uuid, status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProcessNotFound
		}
		return err
	}

	return nil
}
