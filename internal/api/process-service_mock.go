package api

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ProcessSrvcMock struct {
	mock.Mock
}

func (s *ProcessSrvcMock) Submit(ctx context.Context, process *ProcessDTO) (string, error) {
	args := s.Called(ctx, process)
	return args.Get(0).(string), args.Error(1)
}
func (s *ProcessSrvcMock) Get(ctx context.Context, code string, uuid string, page int, pageSize int) (ProcessListDTO, error) {
	args := s.Called(ctx, code, uuid, page, pageSize)
	res := args.Get(0)
	if res != nil {
		return args.Get(0).(ProcessListDTO), args.Error(1)
	}
	return nil, args.Error(1)
}
func (s *ProcessSrvcMock) AssignStatus(ctx context.Context, code string, uuid string, status string, payload Payload) error {
	args := s.Called(ctx, code, uuid, status, payload)
	return args.Error(0)
}
