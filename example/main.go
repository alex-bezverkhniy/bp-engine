package main

import (
	bpengine "bp-engine"
	"flag"
	"log"
	"os"
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
	migrateDbFlag := flag.Bool("migrate", false, "run DB migration scripts")
	serveFlag := flag.Bool("serve", true, "run http server")

	flag.Parse()
	migrateDB := migrateDbFlag != nil && *migrateDbFlag
	serveHTTP := serveFlag != nil && *serveFlag

	cfg, err := bpengine.LoadConfig(confFilePath, env)
	if err != nil {
		log.Fatal("cannot engine config: ", err)
	}

	engine, err := bpengine.New(cfg)
	if err != nil {
		log.Fatal("cannot create new engine: ", err)
	}

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
