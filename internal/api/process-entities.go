package api

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

func (p *Process) toDTO() ProcessDTO {
	var status *ProcessStatusDTO
	if len(p.CurrentStatus.Name) > 0 {
		status = p.CurrentStatus.toDTO()
	}
	return ProcessDTO{
		UUID:          p.UUID,
		Code:          p.Code,
		Payload:       toPayloadDTO(p.Payload),
		CurrentStatus: status,
		Statuses:      p.Statuses.toDTO(),
		CreatedAt:     &p.CreatedAt,
		ChangedAt:     &p.UpdatedAt,
	}
}

func (p *ProcessStatus) toDTO() *ProcessStatusDTO {
	return &ProcessStatusDTO{
		Name:      p.Name,
		Payload:   toPayloadDTO(p.Payload),
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

func (pl ProcessList) toDTO() ProcessListDTO {
	res := ProcessListDTO{}
	for _, p := range pl {
		res = append(res, p.toDTO())
	}
	return res
}

func toPayloadDTO(d datatypes.JSON) Payload {
	val := d.String()
	var payload Payload
	if len(val) > 0 {
		json.Unmarshal([]byte(val), &payload)
	}
	return payload
}
