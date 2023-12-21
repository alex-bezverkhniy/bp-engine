package api

import (
	"bp-engine/internal/model"
	"bp-engine/internal/validators"
	"bytes"
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

func TestSubmit(t *testing.T) {
	defaultUuid := uuid.NewString()
	// ctx := context.Background()
	type args struct {
		reqPayload model.ProcessDTO
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantResp model.ProcessSubmitResponse
		wantErr  *model.ProcessErrorResponse
		mockFunc func(args) *ProcessController
	}{
		{
			name: "fail - 400",
			args: args{
				reqPayload: model.ProcessDTO{
					Code: "requests",
					Payload: model.Payload{
						"sample": "data",
					},
				},
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Submit", mock.Anything, &args.reqPayload).
					Return(defaultUuid, nil)
				return NewProcessController(&service)
			},
			wantCode: http.StatusBadRequest,
			wantErr:  &CannotReadRequestBodyErrResp,
		},
		{
			name: "fail - 500",
			args: args{
				reqPayload: model.ProcessDTO{
					Code: "requests",
					Payload: model.Payload{
						"sample": "data",
					},
				},
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Submit", mock.Anything, &args.reqPayload).
					Return("", errors.New("OMG error"))
				return NewProcessController(&service)
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  &CannotCreateNewProcessErrResp,
		},
		{
			name: "success",
			args: args{
				reqPayload: model.ProcessDTO{
					Code: "requests",
					Payload: model.Payload{
						"sample": "data",
					},
				},
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Submit", mock.Anything, &args.reqPayload).
					Return(defaultUuid, nil)
				return NewProcessController(&service)
			},
			wantCode: http.StatusOK,
			wantResp: model.ProcessSubmitResponse{
				Uuid: defaultUuid,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testApp = fiber.New()
			controller := tt.mockFunc(tt.args)

			testGroup := testApp.Group("/test/")
			controller.SetupRouter(testGroup)
			url := "http://localhost/test/"

			var data []byte
			var err error
			if tt.wantCode == http.StatusBadRequest {
				data = []byte("something bad")
			} else {
				data, err = json.Marshal(tt.args.reqPayload)
				assert.Nil(t, err)
			}

			reqBody := bytes.NewBuffer(data)

			req, err := http.NewRequest("POST", url, reqBody)
			assert.Nil(t, err)

			req.Header.Add("Content-Type", "application/json")

			resp, err := testApp.Test(req)

			assert.Nil(t, err)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			if tt.wantErr != nil {
				var gotResp model.ProcessErrorResponse
				json.Unmarshal(respBody, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, *tt.wantErr, gotResp)

			} else {
				var gotResp model.ProcessSubmitResponse
				json.Unmarshal(respBody, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, tt.wantResp.Uuid, gotResp.Uuid)
			}

		})
	}
}

func TestGetList(t *testing.T) {
	defaultUuid := uuid.NewString()
	// ctx := context.Background()
	type args struct {
		code     string
		page     string
		pageSize string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantResp PaginatedResponse
		wantErr  *model.ProcessErrorResponse
		mockFunc func(args) *ProcessController
	}{
		{
			name: "failed - 400 wrong page value",
			args: args{
				code:     "test",
				page:     "abc",
				pageSize: "5",
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", DEFAULT_PAGE, 5).
					Return(nil, ErrProcessNotFound)
				return NewProcessController(&service)
			},
			wantCode: http.StatusBadRequest,
			wantErr:  &NotSupportedValueForPageHdrErrResp,
		},
		{
			name: "failed - 400 wrong page_size value",
			args: args{
				code:     "test",
				page:     "0",
				pageSize: "abc",
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(nil, ErrProcessNotFound)
				return NewProcessController(&service)
			},
			wantCode: http.StatusBadRequest,
			wantErr:  &NotSupportedValueForPageSizeHdrErrResp,
		},
		{
			name: "failed - 404",
			args: args{
				code:     "test",
				page:     "0",
				pageSize: "0",
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(nil, ErrProcessNotFound)
				return NewProcessController(&service)
			},
			wantCode: http.StatusNotFound,
			wantErr:  &ProcessNotFoundErrResp,
		},
		{
			name: "failed - 500",
			args: args{
				code:     "test",
				page:     "0",
				pageSize: "0",
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(nil, errors.New("OMG error"))
				return NewProcessController(&service)
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  &CannotGetListProcessErrResp,
		},
		{
			name: "success",
			args: args{
				code:     "test",
				page:     "1",
				pageSize: "5",
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", 1, 5).
					Return(model.ProcessListDTO{
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
				Data: model.ProcessListDTO{
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
				page:     "0",
				pageSize: "0",
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, "", DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(model.ProcessListDTO{
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
				Data: model.ProcessListDTO{
					{
						Code: "test",
						UUID: defaultUuid,
					},
				},
			},
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
			req.Header.Add(HEADERNAME_PAGE, tt.args.page)
			req.Header.Add(HEADERNAME_PAGE_SIZE, tt.args.pageSize)

			resp, err := testApp.Test(req)

			assert.Nil(t, err)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			if tt.wantErr != nil {
				var gotResp model.ProcessErrorResponse
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
		wantResp model.ProcessListDTO
		wantErr  *model.ProcessErrorResponse
		mockFunc func(args) *ProcessController
	}{
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
			wantErr:  &ProcessNotFoundErrResp,
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
			wantErr:  &CannotGetProcessErrResp,
		},
		{
			name: "success",
			args: args{
				code: "test",
				uuid: defaultUuid,
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("Get", mock.Anything, args.code, args.uuid, DEFAULT_PAGE, DEFAULT_PAGE_SIZE).
					Return(model.ProcessListDTO{
						{
							Code: "test",
							UUID: defaultUuid,
						},
					}, nil)
				return NewProcessController(&service)
			},
			wantCode: http.StatusOK,
			wantResp: model.ProcessListDTO{
				{
					Code: "test",
					UUID: defaultUuid,
				},
			},
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
				var gotResp model.ProcessErrorResponse
				json.Unmarshal(body, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, *tt.wantErr, gotResp)

			} else {
				var gotResp model.ProcessListDTO
				json.Unmarshal(body, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, tt.wantResp, gotResp)
			}

		})
	}
}

func TestAssignStatus(t *testing.T) {
	defaultUuid := uuid.NewString()
	// ctx := context.Background()
	type args struct {
		code       string
		uuid       string
		status     string
		reqPayload model.ProcessStatusDTO
	}
	tests := []struct {
		name               string
		args               args
		wantCode           int
		simulateBadRequest bool
		wantResp           model.ProcessSubmitResponse
		wantErr            *model.ProcessErrorResponse
		mockFunc           func(args) *ProcessController
	}{
		{
			name: "fail - 400",
			args: args{
				code:   "requests",
				uuid:   defaultUuid,
				status: "done",
				reqPayload: model.ProcessStatusDTO{
					Payload: model.Payload{
						"sample": "data",
					},
				},
			},
			simulateBadRequest: true,
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("AssignStatus", mock.Anything,
					args.code,
					args.uuid,
					args.status,
					&args.reqPayload).
					Return(nil)
				return NewProcessController(&service)
			},
			wantCode: http.StatusBadRequest,
			wantErr:  &CannotReadRequestBodyErrResp,
		},
		{
			name: "fail - 400 - not supported status",
			args: args{
				code:   "requests",
				uuid:   defaultUuid,
				status: "done",
				reqPayload: model.ProcessStatusDTO{
					Payload: model.Payload{
						"sample": "data",
					},
				},
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("AssignStatus", mock.Anything,
					args.code,
					args.uuid,
					args.status,
					args.reqPayload.Payload).
					Return(validators.ErrUnknownStatus)
				return NewProcessController(&service)
			},
			wantCode: http.StatusBadRequest,
			wantErr:  &NotSupportedProcessStatusErrResp,
		},
		{
			name: "fail - 400 - not allowed status",
			args: args{
				code:   "requests",
				uuid:   defaultUuid,
				status: "done",
				reqPayload: model.ProcessStatusDTO{
					Payload: model.Payload{
						"sample": "data",
					},
				},
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("AssignStatus", mock.Anything,
					args.code,
					args.uuid,
					args.status,
					args.reqPayload.Payload).
					Return(validators.ErrNotAllowedStatus)
				return NewProcessController(&service)
			},
			wantCode: http.StatusBadRequest,
			wantErr:  &NotAllowedProcessStatusErrResp,
		},
		{
			name: "fail - 404",
			args: args{
				code:   "requests",
				uuid:   defaultUuid,
				status: "done",
				reqPayload: model.ProcessStatusDTO{
					Payload: model.Payload{
						"sample": "data",
					},
				},
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("AssignStatus", mock.Anything,
					args.code,
					args.uuid,
					args.status,
					args.reqPayload.Payload).
					Return(ErrProcessNotFound)
				return NewProcessController(&service)
			},
			wantCode: http.StatusNotFound,
			wantErr:  &ProcessNotFoundErrResp,
		},
		{
			name: "fail - 500",
			args: args{
				code:   "requests",
				uuid:   defaultUuid,
				status: "done",
				reqPayload: model.ProcessStatusDTO{
					Payload: model.Payload{
						"sample": "data",
					},
				},
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("AssignStatus", mock.Anything,
					args.code,
					args.uuid,
					args.status,
					args.reqPayload.Payload).
					Return(errors.New("OMG error"))
				return NewProcessController(&service)
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  &CannotMoveItIntoNewStatusErrResp,
		},
		{
			name: "success",
			args: args{
				code:   "requests",
				uuid:   defaultUuid,
				status: "done",
				reqPayload: model.ProcessStatusDTO{
					Payload: model.Payload{
						"sample": "data",
					},
				},
			},
			mockFunc: func(args args) *ProcessController {
				service := ProcessSrvcMock{}
				service.On("AssignStatus", mock.Anything,
					args.code,
					args.uuid,
					args.status,
					args.reqPayload.Payload).
					Return(nil)
				return NewProcessController(&service)
			},
			wantCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testApp = fiber.New()
			controller := tt.mockFunc(tt.args)

			testGroup := testApp.Group("/test/")
			controller.SetupRouter(testGroup)
			url := fmt.Sprintf("http://localhost/test/%s/%s/assign/%s", tt.args.code, tt.args.uuid, tt.args.status)

			var data []byte
			var err error
			if tt.simulateBadRequest {
				data = []byte("something bad")
			} else {
				data, err = json.Marshal(tt.args.reqPayload)
				assert.Nil(t, err)
			}

			reqBody := bytes.NewBuffer(data)

			req := httptest.NewRequest("PATCH", url, reqBody)

			req.Header.Add("Content-Type", "application/json")

			resp, err := testApp.Test(req)

			assert.Nil(t, err)
			assert.Equal(t, tt.wantCode, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			if tt.wantErr != nil {
				var gotResp model.ProcessErrorResponse
				json.Unmarshal(respBody, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, *tt.wantErr, gotResp)

			} else {
				var gotResp model.ProcessSubmitResponse
				json.Unmarshal(respBody, &gotResp)
				assert.Nil(t, err)

				assert.Equal(t, tt.wantResp.Uuid, gotResp.Uuid)
			}

		})
	}
}

func Test_getHeaderValue(t *testing.T) {

	type args struct {
		headerKey string
		defVal    any
	}
	tests := []struct {
		name    string
		args    args
		headers map[string][]string
		wantVal any
		wantErr error
	}{
		{
			name: "success",
			args: args{
				headerKey: HEADERNAME_PAGE,
				defVal:    strconv.Itoa(DEFAULT_PAGE),
			},
			headers: map[string][]string{HEADERNAME_PAGE: {"11"}},
			wantVal: "11",
		},
		{
			name: "success - default val",
			args: args{
				headerKey: HEADERNAME_PAGE,
				defVal:    strconv.Itoa(DEFAULT_PAGE),
			},
			wantVal: "1",
		},
		{
			name: "success - default val - int",
			args: args{
				headerKey: HEADERNAME_PAGE,
				defVal:    DEFAULT_PAGE,
			},
			wantVal: 1,
		},
		{
			name: "success - default val - float64",
			args: args{
				headerKey: HEADERNAME_PAGE,
				defVal:    1.5,
			},
			wantVal: 1.5,
		},
		{
			name: "success - val - float64",
			args: args{
				headerKey: HEADERNAME_PAGE,
				defVal:    1.0,
			},
			headers: map[string][]string{HEADERNAME_PAGE: {"1.5"}},
			wantVal: 1.5,
		},
		{
			name: "success - val - int",
			args: args{
				headerKey: HEADERNAME_PAGE,
				defVal:    1,
			},
			headers: map[string][]string{HEADERNAME_PAGE: {"5"}},
			wantVal: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr != nil {
				_, gotErr := getHeaderValue(tt.headers, tt.args.headerKey, tt.args.defVal.(string))
				assert.NotNil(t, gotErr)
				assert.Equal(t, tt.wantErr.Error(), gotErr.Error())
			} else {
				switch tt.wantVal.(type) {
				case int:
					var gotVal int
					gotVal, gotErr := getHeaderValue(tt.headers, tt.args.headerKey, tt.args.defVal.(int))
					assert.Nil(t, gotErr)
					assert.Equal(t, tt.wantVal, gotVal)
				case float64:
					var gotVal float64
					gotVal, gotErr := getHeaderValue(tt.headers, tt.args.headerKey, tt.args.defVal.(float64))
					assert.Nil(t, gotErr)
					assert.Equal(t, tt.wantVal, gotVal)
				default:
					var gotVal string
					gotVal, gotErr := getHeaderValue(tt.headers, tt.args.headerKey, tt.args.defVal.(string))
					assert.Nil(t, gotErr)
					assert.Equal(t, tt.wantVal, gotVal)
				}

			}

		})
	}
}
