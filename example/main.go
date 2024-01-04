package main

import (
	"flag"
	"log"
	"os"

	bpengine "github.com/alex-bezverkhniy/bp-engine"
	"github.com/gofiber/fiber/v2"
)

// @title Example Business Process Engine API
// @version 1.0
// @description This is an example of the Business Process Engine API
// @termsOfService http://swagger.io/terms/
// @contact.name Alex Bezverkhniy
// @contact.email alexander.bezverkhniy@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	env := os.Getenv("ENV")
	confFilePath := os.Getenv("CONFIG_FILE")
	dbURL := os.Getenv("DB_URL")

	migrateDbFlag := flag.Bool("migrate", false, "run DB migration scripts")
	serveFlag := flag.Bool("serve", true, "run http server")

	flag.Parse()
	migrateDB := migrateDbFlag != nil && *migrateDbFlag
	serveHTTP := serveFlag != nil && *serveFlag

	cfg, err := bpengine.LoadConfig(confFilePath, env)
	if err != nil {
		log.Fatal("cannot engine config: ", err)
	}

	cfg.DbUrl = dbURL
	cfg.SwaggerConfig.Title = "Example API"

	engine, err := bpengine.New(cfg)
	if err != nil {
		log.Fatal("cannot create new engine: ", err)
	}

	// Add custom route
	engine.App.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"test": "OK",
		})
	})

	// Static route
	engine.App.Static("/", "./public")

	if migrateDB {
		if err := engine.SetupDB(cfg); err != nil {
			log.Fatal("cannot setup engine db: ", err)
		}
		log.Fatal(engine.RunDBMigration())
	}

	if serveHTTP {
		if err := engine.InitDefault(); err != nil {
			log.Fatal("cannot init default engine: ", err)
		}

		log.Fatal(engine.Listen(":3000"))
	}

}
