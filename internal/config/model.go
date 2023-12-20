package config

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
