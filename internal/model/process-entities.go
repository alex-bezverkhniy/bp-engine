package model

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	ProcessList []Process

	Process struct {
		gorm.Model
		UUID          string
		Code          string
		Payload       datatypes.JSON
		CurrentStatus ProcessStatus
		Statuses      ProcessStatusList
	}

	ProcessStatusList []ProcessStatus

	ProcessStatus struct {
		gorm.Model
		ProcessID uint
		Name      string
		Payload   datatypes.JSON
	}
)

func (p Process) ToDTO() ProcessDTO {
	var status *ProcessStatusDTO
	if len(p.CurrentStatus.Name) > 0 {
		status = p.CurrentStatus.ToDTO()
	}
	return ProcessDTO{
		UUID:          p.UUID,
		Code:          p.Code,
		Payload:       ToDTO(p.Payload),
		CurrentStatus: status,
		Statuses:      p.Statuses.ToDTO(),
		CreatedAt:     &p.CreatedAt,
		ChangedAt:     &p.UpdatedAt,
	}
}

func (p *ProcessStatus) ToDTO() *ProcessStatusDTO {
	return &ProcessStatusDTO{
		Name:      p.Name,
		Payload:   ToDTO(p.Payload),
		CreatedAt: &p.CreatedAt,
	}
}

func (pp ProcessStatusList) ToDTO() ProcessStatusListDTO {
	res := ProcessStatusListDTO{}
	for _, p := range pp {
		res = append(res, *p.ToDTO())
	}

	return res
}

func (pl ProcessList) ToDTO() ProcessListDTO {
	res := ProcessListDTO{}
	for _, p := range pl {
		res = append(res, p.ToDTO())
	}
	return res
}

func ToDTO(d datatypes.JSON) Payload {
	val := d.String()
	var payload Payload
	if len(val) > 0 {
		json.Unmarshal([]byte(val), &payload)
	}
	return payload
}
