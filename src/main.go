package main

import (
	"errors"

	fiber "github.com/gofiber/fiber/v2"
	log "github.com/gofiber/fiber/v2/log"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	db.AutoMigrate(&Process{})
	db.AutoMigrate(&ProcessStatus{})

	// db.Create(&Process{Code: "requests", Metadata: "sample process for multiple requests"})

	// var prc Process
	// db.First(&prc, 1)
	// log.Info("Created new process", prc)

	app := fiber.New()
	app.Use(fiberlogger.New())

	processRepository := NewProcessRepository(db)

	app.Get("/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
		})
	})

	app.Post("/v1/process", func(c *fiber.Ctx) error {
		var process Process
		err := c.BodyParser(&process)
		if err != nil {
			log.Error("cannot read request body ", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{

				"status":  "error",
				"message": "cannot read request body",
			})
		}

		log.Info("create new process: ", process)
		uuid, err := processRepository.Create(c.Context(), &process)

		if err != nil {
			log.Error("cannot create new process ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{

				"status":  "error",
				"message": "cannot create new process",
			})
		}
		process.UUID = uuid
		return c.Status(fiber.StatusOK).JSON(process)
	})

	app.Get("/v1/process/:code/list", func(c *fiber.Ctx) error {
		code := c.Params("code")
		log.Info("get lits of process by code: ", code)
		processesList, err := processRepository.GetByCode(c.Context(), code)

		if err != nil {
			log.Error("cannot get processes list by code ", err)
			if errors.Is(err, ErrNoRecordsFound) {
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
	})

	app.Get("/v1/process/:uuid", func(c *fiber.Ctx) error {
		uuid := c.Params("uuid")
		log.Info("get process by UUID: ", uuid)
		process, err := processRepository.GetByUUID(c.Context(), uuid)

		if err != nil {
			log.Error("cannot get process by UUID ", err)
			if errors.Is(err, ErrNoRecordsFound) {
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
	})

	app.Patch("/v1/process/:uuid/into/:status", func(c *fiber.Ctx) error {
		uuid := c.Params("uuid")
		status := c.Params("status")
		log.Info("get process by uuid: ", uuid)
		log.Info("move it to: ", status)

		err := processRepository.SetStatus(c.Context(), uuid, status)
		if err != nil {
			log.Error("cannot move into new status ", err)
			if errors.Is(err, ErrRecordNotFound) {
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
	})

	log.Fatal(app.Listen(":3000"))
}
