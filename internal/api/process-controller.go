package api

import (
	"bp-engine/internal/model"
	"bp-engine/internal/validators"
	"errors"
	"strconv"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	log "github.com/gofiber/fiber/v2/log"
)

const (
	HEADERNAME_PAGE_SIZE = "X-Page-Size"
	HEADERNAME_PAGE      = "X-Page"
)

type (
	PaginatedResponse struct {
		Data     model.ProcessListDTO `json:"data"`
		Page     int                  `json:"page"`
		PageSize int                  `json:"page_size"`
	}

	ProcessController struct {
		service ProcessService
	}
)

var (
	ProcessNotFoundErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "process not found",
	}

	CannotGetProcessErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "cannot get process by UUID",
	}

	CannotGetListProcessErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "cannot get processes list by code",
	}
	CannotReadRequestBodyErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "cannot read request body",
	}
	CannotCreateNewProcessErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "cannot create new process",
	}
	CannotMoveItIntoNewStatusErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "cannot move into new status",
	}
	NotSupportedValueForPageHdrErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "not supported value for " + HEADERNAME_PAGE,
	}
	NotSupportedValueForPageSizeHdrErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "not supported value for " + HEADERNAME_PAGE_SIZE,
	}
	NotSupportedProcessStatusErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "not supported process status",
	}

	NotAllowedProcessStatusErrResp = model.ProcessErrorResponse{
		Status:  "error",
		Message: "not allowed process status",
	}
)

func NewProcessController(service ProcessService) *ProcessController {
	return &ProcessController{
		service: service,
	}
}

func (pc *ProcessController) SetupRouter(router fiber.Router) {
	router.Post("/", pc.Submit)
	router.Get("/:code/list", pc.GetList)
	router.Get("/:code/:uuid", pc.Get)
	router.Patch("/:code/:uuid/assign/:status", pc.AssignStatus)
}

// @Summary Creates new process
// @Description Submits/Creates new process
// @Tags process
// @Accept application/json
// @Param	request	body	ProcessDTO	true	"ProcessRequest"
// @Produce json
// @Success 200 {object} ProcessSubmitResponse
// @Router /api/v1/process/ [post]
func (pc *ProcessController) Submit(c *fiber.Ctx) error {
	var process model.ProcessDTO
	err := c.BodyParser(&process)
	if err != nil {
		log.Error("cannot read request body ", err)
		return c.Status(fiber.StatusBadRequest).JSON(CannotReadRequestBodyErrResp)
	}

	log.Infof("create new process: %v", process)
	uuid, err := pc.service.Submit(c.Context(), &process)

	if err != nil {
		log.Error("cannot create new process ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(CannotCreateNewProcessErrResp)
	}
	res := model.ProcessSubmitResponse{
		Uuid: uuid,
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

// @Summary Get list of processes
// @Description Get list of processes
// @Tags process
// @Param	code		path	string	true	"Code of Process"
// @Param	X-Page		header	int		false	"Page number"
// @Param	X-Page-Size	header	int		false	"Page size"
// @Produce json
// @Success 200 {object} ProcessListDTO
// @Router /api/v1/process/{code}/list [get]
func (pc *ProcessController) GetList(c *fiber.Ctx) error {
	code := c.Params("code")

	var err error
	var page int
	var pageSize int

	page, err = getHeaderValue[int](c.GetReqHeaders(), HEADERNAME_PAGE, DEFAULT_PAGE)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NotSupportedValueForPageHdrErrResp)
	}
	if page == 0 {
		page = DEFAULT_PAGE
	}

	pageSize, err = getHeaderValue[int](c.GetReqHeaders(), HEADERNAME_PAGE_SIZE, DEFAULT_PAGE_SIZE)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NotSupportedValueForPageSizeHdrErrResp)
	}
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}

	log.Info("get lits of process by code: ", code)
	processesList, err := pc.service.Get(c.Context(), code, "", page, pageSize)

	if err != nil {
		log.Error("cannot get processes list by code ", err)
		if errors.Is(err, ErrProcessNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ProcessNotFoundErrResp)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(CannotGetListProcessErrResp)
	}

	resp := PaginatedResponse{
		Data:     processesList,
		Page:     page,
		PageSize: pageSize,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// @Summary Get process
// @Description Get process by UUID
// @Tags process
// @Param	code	path	string	true	"Code of Process"
// @Param	uuid	path	string	true	"UUID of Process"
// @Produce json
// @Success	200 {object} ProcessListDTO
// @Failed	404 {object} ProcessErrorResponse
// @Failed	500 {object} ProcessErrorResponse
// @Router /api/v1/process/{code}/{uuid} [get]
func (pc *ProcessController) Get(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	code := c.Params("code")
	log.Infof("get process by code: %s and UUID: %s", code, uuid)
	process, err := pc.service.Get(c.Context(), code, uuid, DEFAULT_PAGE, DEFAULT_PAGE_SIZE)

	if err != nil {
		log.Error("cannot get process by UUID ", err)
		if errors.Is(err, ErrProcessNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ProcessNotFoundErrResp)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(CannotGetProcessErrResp)
	}

	return c.Status(fiber.StatusOK).JSON(process)
}

// @Summary Assign the process to the status
// @Description Assign/move the process to the status
// @Tags process
// @Accept application/json
// @Param	code	path	string				true	"Code of Process"
// @Param	uuid	path	string				true	"UUID of Process"
// @Param	status	path	string				true	"Status of Process"
// @Param	request	body	ProcessStatusDTO	true	"ProcessStatus"
// @Produce json
// @Success 204
// @Router /api/v1/process/{code}/{uuid}/assign/{status}	[patch]
func (pc *ProcessController) AssignStatus(c *fiber.Ctx) error {
	code := c.Params("code")
	uuid := c.Params("uuid")
	status := c.Params("status")
	ctx := c.Context()

	var processStatus model.ProcessStatusDTO
	err := c.BodyParser(&processStatus)
	if err != nil {
		log.Error("cannot read request body ", err)
		return c.Status(fiber.StatusBadRequest).JSON(CannotReadRequestBodyErrResp)
	}

	log.Info("get process by uuid: ", uuid)
	log.Info("move it to: ", status)

	err = pc.service.AssignStatus(ctx, code, uuid, status, processStatus.Payload)
	if err != nil {
		log.Error("cannot move into new status ", err)
		if errors.Is(err, ErrProcessNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ProcessNotFoundErrResp)
		}
		if errors.Is(err, validators.ErrUnknownStatus) {
			return c.Status(fiber.StatusBadRequest).JSON(NotSupportedProcessStatusErrResp)
		}
		if errors.Is(err, validators.ErrNotAllowedStatus) {
			return c.Status(fiber.StatusBadRequest).JSON(NotAllowedProcessStatusErrResp)
		}
		if errors.Is(err, validators.ErrPayloadValidation) {
			return c.Status(fiber.StatusBadRequest).JSON(
				model.ProcessErrorResponse{
					Status:  "error",
					Message: strings.ReplaceAll(err.Error(), "\n", ""),
				},
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(CannotMoveItIntoNewStatusErrResp)
	}

	c.Status(fiber.StatusNoContent)
	return nil

}

func getHeaderValue[T string | int | float64](headers map[string][]string, key string, defaultVal T) (T, error) {
	var err error
	var val any
	if headers != nil {
		if len(headers[key]) > 0 {
			strVal := headers[key][0]
			var v any = new(T)
			switch v.(type) {
			case *int:
				val, err = strconv.Atoi(strVal)
			case *float64:
				val, err = strconv.ParseFloat(strVal, 64)
			default: //string
				val = strVal
			}
			return val.(T), err
		}
		return defaultVal, err
	}
	return defaultVal, err
}
