package api

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

type (
	GraphNode struct {
		Name         string
		NextStatuses []GraphNode
	}
	ProcessModel struct {
		Name       string
		GraphModel GraphNode
	}

	ProcessListDTO []ProcessDTO
	ProcessDTO     struct {
		UUID          string               `json:"uuid"`
		Code          string               `json:"code"`
		Metadata      Metadata             `json:"metadata,omitempty"`
		CurrentStatus *ProcessStatusDTO    `json:"current_status,omitempty"`
		Statuses      ProcessStatusListDTO `json:"statuses,omitempty"`
		CreatedAt     time.Time            `json:"created_at"`
		ChangedAt     time.Time            `json:"changed_at"`
	}

	ProcessStatusListDTO []ProcessStatusDTO

	Metadata map[string]interface{}

	ProcessStatusDTO struct {
		Name      string    `json:"name,omitempty"`
		Metadata  Metadata  `json:"metadata,omitempty"`
		CreatedAt time.Time `json:"created_at"`
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
