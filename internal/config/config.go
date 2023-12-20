package config

import (
	"encoding/json"
	"errors"
)

type (
	FileReader interface {
		ReadFile(filePath string) ([]byte, error)
	}

	Config struct {
		filePath      string          `json:"-"`
		ProcessConfig []ProcessConfig `json:"processes"`
	}
)

const DEFAULT_CONFIG_FILEPATH = "./config.json"

var ErrConfigFileIsEmpty = errors.New("Config file is empty")

func (c *Config) LoadConfig(fileReader FileReader) (*Config, error) {
	var conf Config

	content, err := fileReader.ReadFile(c.filePath)
	if err != nil {
		return nil, err
	}

	if len(content) <= 0 {
		return nil, ErrConfigFileIsEmpty
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
