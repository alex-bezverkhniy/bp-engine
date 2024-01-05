package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gofiber/contrib/swagger"
)

const (
	DEFAULT_PAGE_SIZE = 10
	DEFAULT_PAGE      = 1
)

type (
	FileReader interface {
		ReadFile(filePath string) ([]byte, error)
	}
	osFileReader  struct{}
	ConfigBuilder struct {
		filePath string `json:"-"`
	}
	Config struct {
		Env           string            `json:"-"`
		DbEngine      string            `json:"db_engine,omitempty"`
		DbUrl         string            `json:"db_url"`
		ProcessConfig ProcessConfigList `json:"processes"`
		SwaggerConfig swagger.Config    `json:"swagger_config:omitempty"`
	}
)

const DEFAULT_CONFIG_FILEPATH = "./config.json"

var ErrConfigFileIsEmpty = errors.New("Config file is empty")

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		filePath: DEFAULT_CONFIG_FILEPATH,
	}
}

func (cb *ConfigBuilder) WithConfigFile(filePath string) *ConfigBuilder {
	cb.filePath = filePath
	return cb
}

func (cb *ConfigBuilder) LoadConfig() (*Config, error) {
	fileReader := osFileReader{}
	return cb.LoadFromFile(&fileReader)
}

func (cb *ConfigBuilder) LoadFromFile(fileReader FileReader) (*Config, error) {
	var conf Config

	content, err := fileReader.ReadFile(cb.filePath)
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

	// check and set default status type
	for _, pc := range conf.ProcessConfig {
		for _, st := range pc.Statuses {
			if len(st.Type) == 0 {
				st.Type = GenericStatus
			}
		}
	}

	return &conf, nil
}

func (osfr *osFileReader) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
