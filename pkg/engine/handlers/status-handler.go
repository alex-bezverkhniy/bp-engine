package handlers

import (
	"context"

	"github.com/alex-bezverkhniy/bp-engine/internal/model"
)

type (
	StatusHandlers map[string]StatusHandler

	StatusHandler interface {
		Process(ctx context.Context, code string, uuid string, status string, payload model.Payload) error
	}
)
