package validators

import (
	"github.com/alex-bezverkhniy/bp-engine/internal/model"

	"github.com/stretchr/testify/mock"
)

type ValidatorMocked struct {
	mock.Mock
}

func (vm *ValidatorMocked) Validate(process model.ProcessDTO, newStatus model.ProcessStatusDTO) error {
	args := vm.Called(process, newStatus)
	return args.Error(0)
}

func (vm *ValidatorMocked) CompileJsonSchema() error {
	args := vm.Called()
	return args.Error(0)
}
