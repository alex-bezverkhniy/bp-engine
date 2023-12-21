package validators

import (
	"bp-engine/internal/config"
	"bp-engine/internal/model"
	"errors"
)

type (
	Validator interface {
		Validate(process model.ProcessDTO, newStatus model.ProcessStatusDTO) error
	}

	BasicValidator struct {
		conf config.ProcessConfigList
	}
)

var ErrUnknownStatus = errors.New("unknown status")
var ErrNotAllowedStatus = errors.New("not allowed status")

func NewBasicValidator(conf []config.ProcessConfig) Validator {
	return &BasicValidator{
		conf: conf,
	}
}

func (bv *BasicValidator) Validate(process model.ProcessDTO, newStatus model.ProcessStatusDTO) error {
	// Check if status defined
	_, err := bv.conf.GetStatusConfig(process.Code, newStatus.Name)
	if err != nil {
		return ErrUnknownStatus
	}

	// Check current status config
	currentStatusCfg, err := bv.conf.GetStatusConfig(process.Code, process.CurrentStatus.Name)
	if err != nil {
		return ErrUnknownStatus
	}
	for _, s := range currentStatusCfg.Next {
		if s == newStatus.Name {
			return nil
		}
	}

	//TODO: Check payload

	return ErrNotAllowedStatus
}
