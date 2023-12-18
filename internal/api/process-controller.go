package api

import (
	"errors"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	log "github.com/gofiber/fiber/v2/log"
)

const (
	HEADERNAME_PAGE_SIZE = "X-Page-Size"
	HEADERNAME_PAGE      = "X-Page"
)

type (
	PaginatedResponse struct {
		Data     ProcessListDTO `json:"data"`
		Page     int            `json:"page"`
		PageSize int            `json:"page_size"`
	}

	ProcessController struct {
		service ProcessService
	}
)

var (
	ProcessNotFoundErrResp = ProcessErrorResponse{
		Status:  "error",
		Message: "process not found",
	}

	CannotGetProcessErrResp = ProcessErrorResponse{
		Status:  "error",
		Message: "cannot get process by UUID",
	}

	CannotGetListProcessErrResp = ProcessErrorResponse{
		Status:  "error",
		Message: "cannot get processes list by code",
	}
	CannotReadRequestBodyErrResp = ProcessErrorResponse{
		Status:  "error",
		Message: "cannot read request body",
	}
	CannotCreateNewProcessErrResp = ProcessErrorResponse{
		Status:  "error",
		Message: "cannot create new process",
	}
	CannotMoveItIntoNewStatusErrResp = ProcessErrorResponse{
		Status:  "error",
		Message: "cannot move into new status",
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
	var process ProcessDTO
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
	res := ProcessSubmitResponse{
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

	page, err = getHeaderValue[int](c, HEADERNAME_PAGE, DEFAULT_PAGE)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "not supported value for " + HEADERNAME_PAGE,
		})
	}
	if page == 0 {
		page = DEFAULT_PAGE
	}

	pageSize, err = getHeaderValue[int](c, HEADERNAME_PAGE_SIZE, DEFAULT_PAGE_SIZE)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "not supported value for " + HEADERNAME_PAGE_SIZE,
		})
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

	var processStatus ProcessStatusDTO
	err := c.BodyParser(&processStatus)
	if err != nil {
		log.Error("cannot read request body ", err)
		return c.Status(fiber.StatusBadRequest).JSON(CannotReadRequestBodyErrResp)
	}

	log.Info("get process by uuid: ", uuid)
	log.Info("move it to: ", status)

	err = pc.service.AssignStatus(c.Context(), code, uuid, status, processStatus.Payload)
	if err != nil {
		log.Error("cannot move into new status ", err)
		if errors.Is(err, ErrProcessNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ProcessNotFoundErrResp)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(CannotMoveItIntoNewStatusErrResp)
	}

	c.Status(fiber.StatusNoContent)
	return nil

}

func getHeaderValue[T string | int | float64](c *fiber.Ctx, key string, defaultVal T) (T, error) {
	headers := c.GetReqHeaders()
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
				val = &strVal
			}
		}
		return val.(T), err
	}
	return defaultVal, err
}
