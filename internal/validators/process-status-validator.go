package validators

import (
	"bp-engine/internal/model"
)

type (
	Validator interface {
		Validate(status model.ProcessStatusDTO) error
	}
)
