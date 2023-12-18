package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	defaultUuid := uuid.NewString()
	// ctx := context.Background()
	type args struct {
		code string
		uuid string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantResp ProcessListDTO
		wantErr  *ProcessErrorResponse
		mockFunc func(args) *ProcessController
	}{
		{
			name: "success",
			args: args{
				code: "test",
				uuid: defaultUuid,
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, args.uuid, DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(ProcessListDTO{
						{
							Code: "test",
							UUID: defaultUuid,
						},
					}, nil)
				return NewProcessController(&service)
			},
			wantCode: http.StatusOK,
			wantResp: ProcessListDTO{
				{
					Code: "test",
					UUID: defaultUuid,
				},
			},
		},
		{
			name: "failed - 404",
			args: args{
				code: "test",
				uuid: defaultUuid,
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, args.uuid, DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(nil, ErrProcessNotFound)
				return NewProcessController(&service)
			},
			wantCode: http.StatusNotFound,
			wantErr:  &ProcessNotFoundResp,
		},
		{
			name: "failed - 500",
			args: args{
				code: "test",
				uuid: defaultUuid,
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, args.uuid, DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(nil, errors.New("odd error"))
				return NewProcessController(&service)
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  &CannotGetProcessResp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testApp = fiber.New()
			controller := tt.mockFunc(tt.args)

			testGroup := testApp.Group("/test/")
			controller.SetupRouter(testGroup)
			url := fmt.Sprintf("http://localhost/test/%s/%s", tt.args.code, tt.args.uuid)
			req := httptest.NewRequest("GET", url, nil)

			resp, err := testApp.Test(req)

			assert.Nil(t, err)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			if tt.wantErr != nil {
				var gotResp ProcessErrorResponse
				json.Unmarshal(body, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, *tt.wantErr, gotResp)

			} else {
				var gotResp ProcessListDTO
				json.Unmarshal(body, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, tt.wantResp, gotResp)
			}

		})
	}
}

func TestGetList(t *testing.T) {
	defaultUuid := uuid.NewString()
	// ctx := context.Background()
	type args struct {
		code     string
		page     int
		pageSize int
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantResp PaginatedResponse
		wantErr  *ProcessErrorResponse
		mockFunc func(args) *ProcessController
	}{
		{
			name: "success",
			args: args{
				code:     "test",
				page:     1,
				pageSize: 5,
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", args.page, args.pageSize).
					Return(ProcessListDTO{
						{
							Code: "test",
							UUID: defaultUuid,
						},
					}, nil)
				return NewProcessController(&service)
			},
			wantCode: http.StatusOK,
			wantResp: PaginatedResponse{
				Page:     1,
				PageSize: 5,
				Data: ProcessListDTO{
					{
						Code: "test",
						UUID: defaultUuid,
					},
				},
			},
		},
		{
			name: "success - default page and pageSize",
			args: args{
				code:     "test",
				page:     0,
				pageSize: 0,
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(ProcessListDTO{
						{
							Code: "test",
							UUID: defaultUuid,
						},
					}, nil)
				return NewProcessController(&service)
			},
			wantCode: http.StatusOK,
			wantResp: PaginatedResponse{
				Page:     DEFAULT_PAGE,
				PageSize: DEFAULT_PAGE_SIZE,
				Data: ProcessListDTO{
					{
						Code: "test",
						UUID: defaultUuid,
					},
				},
			},
		},
		{
			name: "failed - 404",
			args: args{
				code: "test",
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(nil, ErrProcessNotFound)
				return NewProcessController(&service)
			},
			wantCode: http.StatusNotFound,
			wantErr:  &ProcessNotFoundResp,
		},
		{
			name: "failed - 500",
			args: args{
				code: "test",
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(nil, errors.New("OMG error"))
				return NewProcessController(&service)
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  &CannotGetListProcessResp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testApp = fiber.New()
			controller := tt.mockFunc(tt.args)

			testGroup := testApp.Group("/test/")
			controller.SetupRouter(testGroup)
			url := fmt.Sprintf("http://localhost/test/%s/list", tt.args.code)
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Add(HEADERNAME_PAGE, strconv.Itoa(tt.args.page))
			req.Header.Add(HEADERNAME_PAGE_SIZE, strconv.Itoa(tt.args.pageSize))

			resp, err := testApp.Test(req)

			assert.Nil(t, err)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			if tt.wantErr != nil {
				var gotResp ProcessErrorResponse
				json.Unmarshal(body, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, *tt.wantErr, gotResp)

			} else {
				var gotResp PaginatedResponse
				json.Unmarshal(body, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, tt.wantResp.Page, gotResp.Page)
				assert.Equal(t, tt.wantResp.PageSize, gotResp.PageSize)
				assert.Equal(t, tt.wantResp.Data, gotResp.Data)

			}

		})
	}
}
