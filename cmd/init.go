package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/knadh/koanf/parsers/toml"
	file "github.com/knadh/koanf/providers/file"
	posflag "github.com/knadh/koanf/providers/posflag"

	_ "github.com/go-jet/jet/v2/postgres"
	"github.com/knadh/koanf/v2"
	echo "github.com/labstack/echo/v4"
	_ "github.com/sarthakjdev/wapikit/.db-generated/wapikit/public/model"
	_ "github.com/sarthakjdev/wapikit/.db-generated/wapikit/public/table"
	flag "github.com/spf13/pflag"
)

func initFlags() {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		// ! TODO: trigger the --help command here
		logger.Info(f.FlagUsages())
		os.Exit(0)
	}

	// Register the command line flags.
	f.StringSlice("config", []string{"config-dev.toml"},
		"path to one or more config files (will be merged in order)")
	f.Bool("install", false, "setup database (first time)")
	f.Bool("idempotent", false, "make --install run only if the database isn't already setup")
	f.Bool("upgrade", false, "upgrade database to the current version")
	f.Bool("version", false, "show current version of the build")
	f.Bool("yes", false, "assume 'yes' to prompts during --install/upgrade")
	f.Bool("db-migrate", false, "apply database migrations")
	f.Bool("db-apply", false, "apply database migrations")

	if err := f.Parse(os.Args[1:]); err != nil {
		logger.Error("error loading flags: %v", err)
	}

	if err := koa.Load(posflag.Provider(f, ".", koa), nil); err != nil {
		logger.Error("error loading config: %v", err)
	}
}

func initDb() {
	// check for the database connection
	// check if the database if the database is already setup
	// if not, then setup the database

}

func handleMigrationApply() {

	// ask for the confirmation here first

}

func checkDbState() {

}

func initMigrate() {

}

func initFs() {

}

func loadConfigFiles(filePaths []string, koa *koanf.Koanf) {
	for _, filePath := range filePaths {
		logger.Info("reading config: %s", filePath)
		if err := koa.Load(file.Provider(filePath), toml.Parser()); err != nil {
			if os.IsNotExist(err) {
				logger.Error("config file not found. If there isn't one yet, run --new-config to generate one.")
			}
			logger.Error("error loading config from file: %v.", err)
		}
	}

}

func installApp() {
	// init the database
	// init the filesystem
	// init the config files
	// apply database migrations

}

// initHTTPServer sets up and runs the app's main HTTP server and blocks forever.
func initHTTPServer(app *App) *echo.Echo {
	app.logger.Info("initializing HTTP server")
	var server = echo.New()
	logger := app.logger
	server.HideBanner = true

	// Register app (*App) to be injected into all HTTP handlers.
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			fmt.Println("injecting app into context")
			c.Set("app", &app)
			return next(c)
		}
	})

	// Register all HTTP handlers.
	mountHandlers(server, app)

	// Start the server.
	func() {
		logger.Info("starting HTTP server on localhost:5000")
		if err := server.Start("localhost:5000"); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println("HTTP server shut down")
			} else {
				logger.Error("error starting HTTP server: %v", err)
			}
		}

	}()

	return server
}
