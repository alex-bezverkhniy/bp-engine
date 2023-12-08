package main

import (
	"bp-engine/internal/api"

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

	db.AutoMigrate(&api.Process{})
	db.AutoMigrate(&api.ProcessStatus{})

	app := fiber.New()
	app.Use(fiberlogger.New())

	processRepository := api.NewProcessRepository(db)
	processService := api.NewProcessService(processRepository)
	processController := api.NewProcessController(processService)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", health)
	processController.SetupRouter(v1.Group("/process"))

	log.Fatal(app.Listen(":3000"))
}

func health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "OK",
	})
}
