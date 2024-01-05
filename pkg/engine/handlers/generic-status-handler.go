package handlers

import (
	"context"
	"errors"

	"github.com/alex-bezverkhniy/bp-engine/internal/model"
	"github.com/alex-bezverkhniy/bp-engine/internal/repositories"
	"github.com/alex-bezverkhniy/bp-engine/pkg/engine"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GenericStatusHandler struct {
	repo repositories.ProcessRepository
}

func NewGenericStatusHandler(repo repositories.ProcessRepository) *GenericStatusHandler {
	return &GenericStatusHandler{
		repo: repo,
	}
}

func (gs *GenericStatusHandler) Process(ctx context.Context, code string, uuid string, status string, payload model.Payload) error {

	err := gs.repo.SetStatus(ctx, code, uuid, status, datatypes.JSON(payload.ToBytes()))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return engine.ErrProcessNotFound
		}
		return err
	}

	return nil
}
