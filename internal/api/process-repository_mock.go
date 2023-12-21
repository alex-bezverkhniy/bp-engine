package api

import (
	"bp-engine/internal/model"
	"context"

	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

type ProcessRepoMock struct {
	mock.Mock
}

func (r *ProcessRepoMock) GetByUUID(ctx context.Context, code string, uuid string) (*model.Process, error) {
	args := r.Called(ctx, code, uuid)
	return args.Get(0).(*model.Process), args.Error(1)
}
func (r *ProcessRepoMock) Create(ctx context.Context, process *model.Process) (string, error) {
	args := r.Called(ctx, process)
	return args.Get(0).(string), args.Error(1)
}
func (r *ProcessRepoMock) GetByCode(ctx context.Context, code string, page int, pageSize int) ([]model.Process, error) {
	args := r.Called(ctx, code, page, pageSize)
	return args.Get(0).([]model.Process), args.Error(1)
}
func (r *ProcessRepoMock) SetStatus(ctx context.Context, code string, uuid string, status string, payload datatypes.JSON) error {
	args := r.Called(ctx, code, uuid, status, payload)
	return args.Error(0)
}
