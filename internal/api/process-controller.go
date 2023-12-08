package api

import (
	"errors"

	fiber "github.com/gofiber/fiber/v2"
	log "github.com/gofiber/fiber/v2/log"
)

type ProcessController struct {
	service ProcessService
}

func NewProcessController(service ProcessService) *ProcessController {
	return &ProcessController{
		service: service,
	}
}

func (pc *ProcessController) SetupRouter(router fiber.Router) {
	router.Post("/", pc.Submit)
	router.Get("/:code/list", pc.GetLists)
	router.Get("/:code/:uuid", pc.Get)
	router.Patch("/:code/:uuid/assign/:status", pc.AssignStatus)
}

func (pc *ProcessController) Submit(c *fiber.Ctx) error {
	var process ProcessDTO
	err := c.BodyParser(&process)
	if err != nil {
		log.Error("cannot read request body ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{

			"status":  "error",
			"message": "cannot read request body",
		})
	}

	log.Infof("create new process: %v", process)
	uuid, err := pc.service.Submit(c.Context(), &process)

	if err != nil {
		log.Error("cannot create new process ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{

			"status":  "error",
			"message": "cannot create new process",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"uuid": uuid,
	})
}

func (pc *ProcessController) GetLists(c *fiber.Ctx) error {
	code := c.Params("code")
	log.Info("get lits of process by code: ", code)
	processesList, err := pc.service.Get(c.Context(), code, "")

	if err != nil {
		log.Error("cannot get processes list by code ", err)
		if errors.Is(err, ErrNoProcessesFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{

				"status":  "error",
				"message": "no processes found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{

			"status":  "error",
			"message": "cannot get processes list by code",
		})
	}

	return c.Status(fiber.StatusOK).JSON(processesList)
}

func (pc *ProcessController) Get(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	code := c.Params("code")
	log.Infof("get process by code: %s and UUID: %s", code, uuid)
	process, err := pc.service.Get(c.Context(), code, uuid)

	if err != nil {
		log.Error("cannot get process by UUID ", err)
		if errors.Is(err, ErrNoProcessesFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{

				"status":  "error",
				"message": "no process found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{

			"status":  "error",
			"message": "cannot get process by UUID",
		})
	}

	return c.Status(fiber.StatusOK).JSON(process)
}

func (pc *ProcessController) AssignStatus(c *fiber.Ctx) error {
	code := c.Params("code")
	uuid := c.Params("uuid")
	status := c.Params("status")

	var processStatus ProcessStatusDTO
	err := c.BodyParser(&processStatus)
	if err != nil {
		log.Error("cannot read request body ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{

			"status":  "error",
			"message": "cannot read request body",
		})
	}

	log.Info("get process by uuid: ", uuid)
	log.Info("move it to: ", status)

	err = pc.service.AssignStatus(c.Context(), code, uuid, status, processStatus.Metadata)
	if err != nil {
		log.Error("cannot move into new status ", err)
		if errors.Is(err, ErrProcessNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{

				"status":  "error",
				"message": "process not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "cannot move into new status",
		})
	}

	c.Status(fiber.StatusNoContent)
	return nil

}
