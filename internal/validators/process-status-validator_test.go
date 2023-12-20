package validators

import (
	"bp-engine/internal/config"
	"bp-engine/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	defaultProcessCOnfig := &config.ProcessConfig{
		Name: "requests",
		Statuses: []config.StatusConfig{
			{
				Name: "open",
				Next: []string{"in_progress", "rejected"},
			},
			{
				Name: "in_progress",
				Next: []string{"open", "rejected", "done"},
			},
			{
				Name: "rejected",
			},
			{
				Name: "done",
			},
		},
	}
	tests := []struct {
		name    string
		conf    *config.ProcessConfig
		process model.ProcessDTO
		status  model.ProcessStatusDTO
		wantErr error
	}{
		{
			name: "valid - move to next",
			conf: defaultProcessCOnfig,
			process: model.ProcessDTO{
				Code: "requests",
				CurrentStatus: &model.ProcessStatusDTO{
					Name: "open",
				},
			},
			status: model.ProcessStatusDTO{
				Name: "in_progress",
			},
			wantErr: nil,
		},
		{
			name: "valid - move to prev",
			conf: defaultProcessCOnfig,
			process: model.ProcessDTO{
				Code: "requests",
				CurrentStatus: &model.ProcessStatusDTO{
					Name: "in_progress",
				},
			},
			status: model.ProcessStatusDTO{
				Name: "open",
			},
			wantErr: nil,
		},
		{
			name: "valid - move to - done",
			conf: defaultProcessCOnfig,
			process: model.ProcessDTO{
				Code: "requests",
				CurrentStatus: &model.ProcessStatusDTO{
					Name: "in_progress",
				},
			},
			status: model.ProcessStatusDTO{
				Name: "done",
			},
			wantErr: nil,
		},
		{
			name: "valid - move to - rejected",
			conf: defaultProcessCOnfig,
			process: model.ProcessDTO{
				Code: "requests",
				CurrentStatus: &model.ProcessStatusDTO{
					Name: "in_progress",
				},
			},
			status: model.ProcessStatusDTO{
				Name: "rejected",
			},
			wantErr: nil,
		},
		{
			name: "invalid - not allowed - unknown",
			conf: defaultProcessCOnfig,
			process: model.ProcessDTO{
				Code: "requests",
				CurrentStatus: &model.ProcessStatusDTO{
					Name: "open",
				},
			},
			status: model.ProcessStatusDTO{
				Name: "complete",
			},
			wantErr: ErrUnknownStatus,
		},
		{
			name: "invalid - not allowed",
			conf: defaultProcessCOnfig,
			process: model.ProcessDTO{
				Code: "requests",
				CurrentStatus: &model.ProcessStatusDTO{
					Name: "open",
				},
			},
			status: model.ProcessStatusDTO{
				Name: "done",
			},
			wantErr: ErrNotAllowedStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewBasicValidator(tt.conf)

			gotErr := validator.Validate(tt.process, tt.status)

			if tt.wantErr != nil {
				assert.NotNil(t, gotErr)
				assert.Equal(t, tt.wantErr.Error(), gotErr.Error())
			} else {
				assert.Nil(t, gotErr)
			}
		})
	}
}
