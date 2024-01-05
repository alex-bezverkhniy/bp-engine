package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedFileReader struct {
	mock.Mock
}

func (fr *mockedFileReader) ReadFile(filePath string) ([]byte, error) {
	args := fr.Called(filePath)
	return args.Get(0).([]byte), args.Error(1)
}

func Test_LoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  error
		wantConf *Config
		mockFunc func() FileReader
	}{
		{
			name: "success",
			mockFunc: func() FileReader {
				fr := mockedFileReader{}
				fr.On("ReadFile", mock.Anything).
					Return(
						[]byte(`{
					"db_url": "gorm.db",
					"processes": [
					{
						"name": "requests",
						"statuses": [
							{
								"name": "open",
								"next": [
									"in_progress",
									"rejected"
								]
							},
							{
								"name": "in_progress",
								"next": [
									"open",
									"rejected"
								]
							},
							{
								"name": "rejected"
							}
						]
					}]		}
					`), nil)
				return &fr
			},
			wantConf: &Config{
				DbUrl: "gorm.db",
				ProcessConfig: []ProcessConfig{
					{
						Name: "requests",
						Statuses: []StatusConfig{
							{
								Name: "open",
								Next: []string{"in_progress", "rejected"},
							},
							{
								Name: "in_progress",
								Next: []string{"open", "rejected"},
							},
							{
								Name: "rejected",
							},
						},
					},
				},
			},
		},
		{
			name: "failed",
			mockFunc: func() FileReader {
				fr := mockedFileReader{}
				fr.On("ReadFile", mock.Anything).
					Return(
						[]byte(`
					{						
							{
								"name": "open",
								"next": [
									"in_progress",
									"rejected"
								]
							},
							{
								"name": "in_progress",
								"next": [
									"open",
									"rejected"
								]
							},
							{
								"name": "rejected"
							}
						]
					}					
					`), nil)
				return &fr
			},
			wantConf: nil,
			wantErr:  errors.New("invalid character '{' looking for beginning of object key string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := NewConfigBuilder()
			fr := tt.mockFunc()
			gotConf, gotErr := conf.LoadFromFile(fr)

			if tt.wantErr != nil {
				assert.NotNil(t, gotErr)
				assert.Equal(t, tt.wantErr.Error(), gotErr.Error())
			} else {
				assert.Nil(t, gotErr)
				assert.NotNil(t, gotConf)
				assert.Equal(t, tt.wantConf, gotConf)
			}

		})
	}
}
