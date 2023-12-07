package main

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	Process struct {
		gorm.Model
		UUID          string
		Code          string
		Metadata      datatypes.JSON
		CurrentStatus ProcessStatus
		Statuses      ProcessStatusList
	}

	ProcessStatusList []ProcessStatus

	ProcessStatus struct {
		gorm.Model
		ProcessID uint
		Name      string
		Metadata  datatypes.JSON
	}
)

func (p *Process) toDTO() ProcessDTO {
	var status *ProcessStatusDTO
	if len(p.CurrentStatus.Name) > 0 {
		status = p.CurrentStatus.toDTO()
	}
	return ProcessDTO{
		UUID:          p.UUID,
		Code:          p.Code,
		Metadata:      toMetadataDTO(p.Metadata),
		CurrentStatus: status,
		Statuses:      p.Statuses.toDTO(),
		CreatedAt:     p.CreatedAt,
		ChangedAt:     p.UpdatedAt,
	}
}

func (p *ProcessStatus) toDTO() *ProcessStatusDTO {
	return &ProcessStatusDTO{
		Name:      p.Name,
		Metadata:  toMetadataDTO(p.Metadata),
		CreatedAt: p.CreatedAt,
	}
}

func (pp ProcessStatusList) toDTO() ProcessStatusListDTO {
	res := ProcessStatusListDTO{}
	for _, p := range pp {
		res = append(res, *p.toDTO())
	}

	return res
}

func toMetadataDTO(d datatypes.JSON) Metadata {
	val := d.String()
	var metadata Metadata
	if len(val) > 0 {
		json.Unmarshal([]byte(val), &metadata)
	}
	return metadata
}
