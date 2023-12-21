package model

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/datatypes"
)

type (
	ProcessListDTO []ProcessDTO
	ProcessDTO     struct {
		UUID          string               `json:"uuid,omitempty" example:"23c968a6-5fc5-4e42-8f59-a7f9c0d4999c"`
		Code          string               `json:"code" example:"requests"`
		Payload       Payload              `json:"payload,omitempty"`
		CurrentStatus *ProcessStatusDTO    `json:"current_status,omitempty"`
		Statuses      ProcessStatusListDTO `json:"statuses,omitempty"`
		CreatedAt     *time.Time           `json:"created_at,omitempty" example:"2023-12-08T11:33:55.418484002-06:00"`
		ChangedAt     *time.Time           `json:"changed_at,omitempty" example:"2023-12-10T12:30:55.442484002-06:00"`
	}

	ProcessStatusListDTO []ProcessStatusDTO

	Payload map[string]interface{}

	ProcessStatusDTO struct {
		Name      string     `json:"name,omitempty" example:"created"`
		Payload   Payload    `json:"payload,omitempty"`
		CreatedAt *time.Time `json:"created_at,omitempty" example:"2023-12-08T11:33:55.418484002-06:00"`
	}

	// Submit process response
	// @Description Response with UUID of created process.
	ProcessSubmitResponse struct {
		Uuid string `json:"uuid" example:"23c968a6-5fc5-4e42-8f59-a7f9c0d4999c"`
	}

	// @Description Error message
	ProcessErrorResponse struct {
		Status  string `json:"status" example:"error"`
		Message string `json:"message" example:"no process found"`
	}
)

func (p *ProcessDTO) ToEntity() *Process {
	var statuses ProcessStatusList
	if len(p.Statuses) > 0 {
		statuses = p.Statuses.ToEntity()
	}

	var curentStatus ProcessStatus
	if p.CurrentStatus != nil {
		curentStatus = p.CurrentStatus.ToEntity()
	}

	return &Process{
		UUID:          p.UUID,
		Code:          p.Code,
		Payload:       p.Payload.ToBytes(),
		CurrentStatus: curentStatus,
		Statuses:      statuses,
	}
}

func (p *ProcessStatusDTO) ToEntity() ProcessStatus {
	metadata := datatypes.JSON{}
	metadata.Scan(p.Payload)
	return ProcessStatus{
		Name:    p.Name,
		Payload: metadata,
	}
}

func (pp ProcessStatusListDTO) ToEntity() ProcessStatusList {
	res := ProcessStatusList{}
	for _, p := range pp {
		res = append(res, p.ToEntity())
	}

	return res
}

func (m Payload) ToBytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}

func (p Payload) ToStringKeys(val interface{}) (interface{}, error) {
	var err error
	switch val := val.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			k, ok := k.(string)
			if !ok {
				return nil, errors.New("found non-string key")
			}
			m[k], err = p.ToStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case []interface{}:
		var l = make([]interface{}, len(val))
		for i, v := range val {
			l[i], err = p.ToStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return val, nil
	}
}
