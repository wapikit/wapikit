package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	file "github.com/knadh/koanf/providers/file"
	posflag "github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/stuffbin"

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

func joinFSPaths(root string, paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		// real_path:stuffbin_alias
		f := strings.Split(p, ":")

		out = append(out, path.Join(root, f[0])+":"+f[1])
	}

	return out
}

// initFileSystem initializes the stuffbin FileSystem to provide
// access to bundled static assets to the app.
func initFS(appDir, frontendDir string) stuffbin.FileSystem {
	var (
		// These paths are joined with "." which is appDir.
		appFiles = []string{
			"./config.toml.sample:config.toml.sample",
		}

		// These path are joined with frontend/out dir
		frontendFiles = []string{
			// frontend/out files should be available on the root path following the file path .
			"./:/",
		}

		// ! TODO: add a static dir path if somebody mounts any other static directory here
	)

	// Get the executable's execPath.
	execPath, err := os.Executable()
	if err != nil {
		logger.Error("error getting executable path: %v", err)
	}

	// Load embedded files in the executable.
	hasEmbed := true
	fs, err := stuffbin.UnStuff(execPath)
	logger.Info("loading embedded filesystem %s", fs.List(), nil)
	if err != nil {
		hasEmbed = false
		// Running in local mode. Load local assets into
		// the in-memory stuffbin.FileSystem.
		logger.Info("unable to initialize embedded filesystem (%v). Using local filesystem", err)
		fs, err = stuffbin.NewLocalFS("/")
		if err != nil {
			logger.Error("failed to initialize local file for assets: %v", err)
		}
	}

	// If the embed failed, load app and frontend files from the compile-time paths.
	files := []string{}
	if !hasEmbed {
		files = append(files, joinFSPaths(appDir, appFiles)...)
		files = append(files, joinFSPaths(frontendDir, frontendFiles)...)
	}

	// No additional files to load.
	if len(files) == 0 {
		return fs
	}

	// Load files from disk and overlay into the FS.
	fStatic, err := stuffbin.NewLocalFS("/", files...)
	if err != nil {
		logger.Error("failed reading static files from disk: '%s': %v", err)
	}

	if err := fs.Merge(fStatic); err != nil {
		logger.Error("error merging static files: '%s': %v", err)
	}

	return fs
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
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", &app)
			return next(c)
		}
	})

	// we want to mount the next.js output to "/" , i.e, / -> "index.html" , /about -> "about.html"
	fileServer := app.fs.FileServer()
	server.GET("/*", echo.WrapHandler(fileServer))

	// Mounting all HTTP handlers.
	mountHandlers(server, app)

	// getting th server address from config and falling back to localhost:5000
	serverAddress := koa.String("address")
	if serverAddress == "" {
		serverAddress = "localhost:5000"
	}

	// Start the server.
	func() {
		logger.Info("starting HTTP server on %s", serverAddress, nil) // Add a placeholder value as the final argument
		if err := server.Start(serverAddress); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println("HTTP server shut down")
			} else {
				logger.Error("error starting HTTP server: %v", err)
			}
		}
	}()

	return server
}
