package api

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type (
	ProcessListDTO []ProcessDTO
	ProcessDTO     struct {
		UUID          string               `json:"uuid" example:"23c968a6-5fc5-4e42-8f59-a7f9c0d4999c"`
		Code          string               `json:"code" example:"requests"`
		Metadata      Metadata             `json:"metadata,omitempty"`
		CurrentStatus *ProcessStatusDTO    `json:"current_status,omitempty"`
		Statuses      ProcessStatusListDTO `json:"statuses,omitempty"`
		CreatedAt     time.Time            `json:"created_at" example:"2023-12-08T11:33:55.418484002-06:00"`
		ChangedAt     time.Time            `json:"changed_at" example:"2023-12-10T12:30:55.442484002-06:00"`
	}

	ProcessStatusListDTO []ProcessStatusDTO

	Metadata map[string]interface{}

	ProcessStatusDTO struct {
		Name      string    `json:"name,omitempty" example:"created"`
		Metadata  Metadata  `json:"metadata,omitempty"`
		CreatedAt time.Time `json:"created_at" example:"2023-12-08T11:33:55.418484002-06:00"`
	}

	// Submit process response
	// @Description Response with UUID of created process.
	ProcessSubmitResponse struct {
		Uuid string `json:"uuid" example:"23c968a6-5fc5-4e42-8f59-a7f9c0d4999c"`
	}
)

func (p *ProcessDTO) toEntity() *Process {
	var statuses ProcessStatusList
	if len(p.Statuses) > 0 {
		statuses = p.Statuses.toEntity()
	}

	var curentStatus ProcessStatus
	if p.CurrentStatus != nil {
		curentStatus = p.CurrentStatus.toEntity()
	}

	return &Process{
		UUID:          p.UUID,
		Code:          p.Code,
		Metadata:      p.Metadata.toBytes(),
		CurrentStatus: curentStatus,
		Statuses:      statuses,
	}
}

func (p *ProcessStatusDTO) toEntity() ProcessStatus {
	metadata := datatypes.JSON{}
	metadata.Scan(p.Metadata)
	return ProcessStatus{
		Name:     p.Name,
		Metadata: metadata,
	}
}

func (pp ProcessStatusListDTO) toEntity() ProcessStatusList {
	res := ProcessStatusList{}
	for _, p := range pp {
		res = append(res, p.toEntity())
	}

	return res
}

func (m Metadata) toBytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}
