package main

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type (
	ProcessDTO struct {
		UUID          string               `json:"uuid"`
		Code          string               `json:"code"`
		Metadata      Metadata             `json:"metadata,omitempty"`
		CurrentStatus *ProcessStatusDTO    `json:"current_status,omitempty"`
		Statuses      ProcessStatusListDTO `json:"statuses,omitempty"`
	}

	ProcessStatusListDTO []ProcessStatusDTO

	Metadata map[string]interface{}

	ProcessStatusDTO struct {
		Name     string   `json: "name"`
		Metadata Metadata `json: "metadata"`
	}
)

func (p *ProcessDTO) toEntity() *Process {
	var statuses ProcessStatusList
	if len(p.Statuses) > 0 {
		statuses = p.Statuses.toEntity()
	}

	var curentStatus ProcessStatus
	if len(p.CurrentStatus.Name) > 0 {
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
