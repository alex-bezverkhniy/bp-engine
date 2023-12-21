package config

import "errors"

type (
	ProcessConfig struct {
		Name     string         `json:"name"`
		Statuses []StatusConfig `json:"statuses"`
	}

	ProcessConfigList []ProcessConfig

	StatusConfig struct {
		Name   string   `json:"name"`
		Next   []string `json:"next,omitempty"`
		Schema string   `json:"schema,omitempty"`
	}
)

var ErrStatusConfigNotFound = errors.New("status config not found")

func (pc ProcessConfigList) GetStatusConfig(code, status string) (*StatusConfig, error) {
	for _, p := range pc {
		if p.Name == code {
			for _, sc := range p.Statuses {
				if sc.Name == status {
					return &sc, nil
				}
			}
		}
	}
	return nil, ErrStatusConfigNotFound
}
