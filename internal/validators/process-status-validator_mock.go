package validators

import (
	"bp-engine/internal/model"

	"github.com/stretchr/testify/mock"
)

type ValidatorMocked struct {
	mock.Mock
}

func (vm *ValidatorMocked) Validate(process model.ProcessDTO, newStatus model.ProcessStatusDTO) error {
	args := vm.Called(process, newStatus)
	return args.Error(0)
}
