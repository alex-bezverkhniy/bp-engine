package main

import (
	"gorm.io/gorm"
)

type (
	Process struct {
		gorm.Model
		UUID     string          `json: "uuid"`
		Code     string          `json: "code"`
		Metadata string          `json: "metadata"`
		Statuses []ProcessStatus `json: "statuses"`
	}

	ProcessStatus struct {
		gorm.Model
		ProcessID uint
		Name      string `json: "name"`
	}
)
