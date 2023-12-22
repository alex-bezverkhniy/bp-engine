package main

import (
	"flag"
	"os"

	"github.com/alex-bezverkhniy/bp-engine/internal/api"
	"github.com/alex-bezverkhniy/bp-engine/internal/config"
	"github.com/alex-bezverkhniy/bp-engine/internal/model"
	"github.com/alex-bezverkhniy/bp-engine/internal/validators"

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
	confFilePath := os.Getenv("CONFIG_FILE")
	migrateDbFlag := flag.Bool("migrate", false, "run DB migration scripts")
	serveFlag := flag.Bool("serve", true, "run http server")

	flag.Parse()
	migrateDB := migrateDbFlag != nil && *migrateDbFlag
	serveHTTP := serveFlag != nil && *serveFlag

	// load config
	conf, err := loadConfig(confFilePath, environment)
	if err != nil {
		log.Fatal("cannot load config file. ", err)
	}

	// connect to db
	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open(conf.DbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	if migrateDB {
		log.Info("run DB migration")
		if migrationErr := runDBMigration(db); len(migrationErr) > 0 {
			for _, e := range migrationErr {
				log.Error("DB migration error", e)
			}
			log.Fatal("cannot successfully complete DB migration")
		}
		log.Info("DB migration done!")
		os.Exit(0)
	}

	if serveHTTP {
		validator, err := setupValidator(conf)
		if err != nil {
			log.Fatal("cannot setup JSON Schema validator", err)
		}

		app := setupApp(conf, validator, db)
		log.Fatal(app.Listen(":3000"))
	}

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

func setupApp(conf *config.Config, validator validators.Validator, db *gorm.DB) *fiber.App {
	// Swagger config
	pathToSwaggerFile := "./docs/swagger.json"
	if len(conf.Env) == 0 || conf.Env == "dev" {
		pathToSwaggerFile = "../docs/swagger.json"
	}
	cfg := swagger.Config{
		BasePath: "/",
		FilePath: pathToSwaggerFile,
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}

	// Fiber App init
	app := fiber.New()
	app.Use(fiberlogger.New())
	app.Use(swagger.New(cfg))

	processRepository := api.NewProcessRepository(db)
	processService := api.NewProcessService(processRepository, validator)
	processController := api.NewProcessController(processService)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", Health)
	processController.SetupRouter(v1.Group("/process"))

	return app
}

func loadConfig(filePath, environment string) (*config.Config, error) {
	if len(filePath) == 0 {
		filePath = config.DEFAULT_CONFIG_FILEPATH
	}
	cb := config.NewConfigBuilder().WithConfigFile(filePath)
	cf, err := cb.LoadConfig()
	if err == nil {
		cf.Env = environment
	}
	return cf, err
}

func runDBMigration(db *gorm.DB) []error {
	var migrationErr []error
	if dbErr := db.AutoMigrate(&model.Process{}); dbErr != nil {
		migrationErr = append(migrationErr, dbErr)
	}
	if dbErr := db.AutoMigrate(&model.ProcessStatus{}); dbErr != nil {
		migrationErr = append(migrationErr, dbErr)
	}

	return migrationErr
}

func setupValidator(conf *config.Config) (validators.Validator, error) {
	validator := validators.NewBasicValidator(conf.ProcessConfig)
	err := validator.CompileJsonSchema()

	if err != nil {
		return nil, err
	}
	return validator, nil
}
