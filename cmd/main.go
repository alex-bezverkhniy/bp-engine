package main

import (
	"bp-engine/internal/api"
	"os"

	"github.com/gofiber/contrib/swagger"
	fiber "github.com/gofiber/fiber/v2"
	log "github.com/gofiber/fiber/v2/log"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title Business Process Engine API
// @version 1.0
// @description This is the Business Process Engine API
// @termsOfService http://swagger.io/terms/
// @contact.name Alex Bezverkhniy
// @contact.email alexander.bezverkhniy@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	environment := os.Getenv("ENV")
	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	db.AutoMigrate(&api.Process{})
	db.AutoMigrate(&api.ProcessStatus{})

	pathToSwaggerFile := "./docs/swagger.json"
	if len(environment) == 0 || environment == "dev" {
		pathToSwaggerFile = "../docs/swagger.json"
	}
	cfg := swagger.Config{
		BasePath: "/",
		FilePath: pathToSwaggerFile,
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}

	app := fiber.New()
	app.Use(fiberlogger.New())
	app.Use(swagger.New(cfg))

	processRepository := api.NewProcessRepository(db)
	processService := api.NewProcessService(processRepository)
	processController := api.NewProcessController(processService)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", Health)
	processController.SetupRouter(v1.Group("/process"))

	log.Fatal(app.Listen(":3000"))
}

// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/Health [get]
func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "OK",
	})
}
