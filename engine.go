package bpengine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alex-bezverkhniy/bp-engine/internal/api"
	"github.com/alex-bezverkhniy/bp-engine/internal/config"
	"github.com/alex-bezverkhniy/bp-engine/internal/model"
	"github.com/alex-bezverkhniy/bp-engine/internal/validators"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrAppIsNotInitialized       = errors.New("fiber app is not initialized")
	ErrDbIsNotInitialized        = errors.New("db is not initialized")
	ErrValidatorIsNotInitialized = errors.New("validator is not initialized")
)

type Engine struct {
	config    config.Config
	app       *fiber.App
	db        *gorm.DB
	validator validators.Validator
}

func New(config config.Config) (*Engine, error) {
	engine := &Engine{
		config: config,
	}

	// Fiber App init
	engine.app = fiber.New()
	engine.app.Use(logger.New())

	return engine, nil
}

func (e *Engine) InitDefault() error {
	// Setup Swagger with default config
	if err := e.SetupSwagger(""); err != nil {
		return err
	}

	// Setup DB
	if err := e.SetupDB(e.config); err != nil {
		return err
	}

	// Setup Validator
	if err := e.SetupValidator(e.config.ProcessConfig); err != nil {
		return err
	}

	// Setup Router
	if err := e.SetupApi(); err != nil {
		return err
	}
	return nil
}

func (e *Engine) Listen(addr string) error {
	if err := e.checkEngineInitialized(); err != nil {
		return err
	}
	return e.app.Listen(addr)
}

func (e *Engine) ListenTLS(addr, certFile, keyFile string) error {
	if err := e.checkEngineInitialized(); err != nil {
		return err
	}
	return e.app.ListenTLS(addr, certFile, keyFile)
}

func (e *Engine) SetupSwagger(pathToSwaggerFile string) error {
	if e.app == nil {
		return ErrAppIsNotInitialized
	}
	// default swagger config
	if len(pathToSwaggerFile) <= 0 {
		// Swagger config
		pathToSwaggerFile = "./internal/docs/swagger.json"
		if len(e.config.Env) == 0 || e.config.Env == "dev" {
			pathToSwaggerFile = "../internal/docs/swagger.json"
		}
	}

	if len(e.config.SwaggerConfig.Path) <= 0 {
		e.config.SwaggerConfig.Path = "swagger"
	}

	e.config.SwaggerConfig.FilePath = pathToSwaggerFile
	e.app.Use(swagger.New(e.config.SwaggerConfig))

	return nil
}

func (e *Engine) SetupDB(cfg config.Config) error {
	var err error

	// github.com/mattn/go-sqlite3
	if len(cfg.DbEngine) == 0 || cfg.DbEngine == "sqlite3" {
		// check file exist
		if _, err := os.Stat(cfg.DbUrl); errors.Is(err, os.ErrNotExist) {
			// set default
			// get current dir
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				return err
			}

			cfg.DbUrl = fmt.Sprintf("%s%c%s", dir, os.PathSeparator, "github.com/alex-bezverkhniy/bp-engine.db")
		}

		e.db, err = gorm.Open(sqlite.Open(cfg.DbUrl), &gorm.Config{})
	}

	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) RunDBMigration() error {
	if e.db == nil {
		return ErrDbIsNotInitialized
	}

	var migrationErr []error
	if dbErr := e.db.AutoMigrate(&model.Process{}); dbErr != nil {
		migrationErr = append(migrationErr, dbErr)
	}
	if dbErr := e.db.AutoMigrate(&model.ProcessStatus{}); dbErr != nil {
		migrationErr = append(migrationErr, dbErr)
	}

	if len(migrationErr) > 0 {
		return errors.Join(migrationErr...)
	}

	return nil
}

func (e *Engine) SetupValidator(cfg config.ProcessConfigList) error {

	// Override default config if needed
	if len(cfg) != 0 {
		e.config.ProcessConfig = cfg
	}

	e.validator = validators.NewBasicValidator(e.config.ProcessConfig)
	err := e.validator.CompileJsonSchema()

	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) SetValidator(customValidator validators.Validator) {
	e.validator = customValidator
}

func (e *Engine) SetupApi() error {
	if err := e.checkEngineInitialized(); err != nil {
		return err
	}

	processRepository := api.NewProcessRepository(e.db)
	processService := api.NewProcessService(processRepository, e.validator)
	processController := api.NewProcessController(processService)

	api := e.app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", Health)
	processController.SetupRouter(v1.Group("/process"))

	return nil
}

func (e *Engine) checkEngineInitialized() error {
	if e.app == nil {
		return ErrAppIsNotInitialized
	}
	if e.db == nil {
		return ErrDbIsNotInitialized
	}
	if e.validator == nil {
		return ErrValidatorIsNotInitialized
	}

	return nil
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

func LoadConfig(filePath, environment string) (config.Config, error) {
	if len(filePath) == 0 {
		filePath = config.DEFAULT_CONFIG_FILEPATH
	}
	cb := config.NewConfigBuilder().WithConfigFile(filePath)
	cf, err := cb.LoadConfig()
	if err == nil && cf != nil {
		cf.Env = environment
		return *cf, err
	}

	return config.Config{}, err
}
