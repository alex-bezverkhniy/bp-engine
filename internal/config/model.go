package config

import "errors"

type (
	ProcessConfig struct {
		Name     string         `json:"name"`
		Statuses []StatusConfig `json:"statuses"`
	}

	StatusConfig struct {
		Name string   `json:"name"`
		Next []string `json:"next,omitempty"`
	}
)

var ErrStatusConfigNotFound = errors.New("status config not found")

func (pc *ProcessConfig) GetStatusConfig(code string) (*StatusConfig, error) {
	for _, sc := range pc.Statuses {
		if sc.Name == code {
			return &sc, nil
		}
	}
	return nil, ErrStatusConfigNotFound
}
