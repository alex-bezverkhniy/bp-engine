package main

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	Process struct {
		gorm.Model
		UUID          string         `json: "uuid"`
		Code          string         `json: "code"`
		Metadata      datatypes.JSON `json: "metadata"`
		CurrentStatus ProcessStatus
		Statuses      []ProcessStatus `json: "statuses"`
	}

	ProcessStatus struct {
		gorm.Model
		ProcessID uint
		Name      string         `json: "name"`
		Metadata  datatypes.JSON `json: "metadata"`
	}
)
