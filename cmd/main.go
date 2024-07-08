package main

import (
	"log/slog"
	"os"
	"sync"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
	api "github.com/sarthakjdev/wapikit/api/cmd"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
	websocket_server "github.com/sarthakjdev/wapikit/websocket-server"
)

// because this will be a single binary, we will be providing the flags here
// 1. --install to install the setup the app, but it will be idempotent
// 2. --migrate to apply the migration to the database
// 3. --config to setup the config files
// 4. --version to check the version of the application running
// 5. --help to check the

var (
	// Global variables
	logger      = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	koa         = koanf.New(".")
	fs          stuffbin.FileSystem
	appDir      string = "."
	frontendDir string = "frontend/out"
)

func init() {
	initFlags()
	loadConfigFiles(koa.Strings("config"), koa)

	if koa.Bool("version") {
		logger.Info("current version of the application")
	}

	// here appDir is for config file packing, frontendDir is for the frontend built output and static dir is any other static files and the public
	fs = initFS(appDir, frontendDir)

	if koa.Bool("install") {
		logger.Info("Installing the application")
		// ! should be idempotent
		installApp(koa.String("last_version"), database.GetDbInstance(), fs, koa.Bool("yes"), koa.Bool("idempotent"))
		os.Exit(0)
	}

	if koa.Bool("upgrade") {
		logger.Info("Upgrading the application")
		// ! should not upgrade without asking for thr permission, because database migration can be destructive
		// upgrade handler
	}
	// do nothing
	// ** NOTE: if no flag is provided, then let the app move to the main function and start the server
}

func main() {
	logger.Info("Starting the application")

	redisUrl := koa.String("app.redis_url")

	if redisUrl == "" {
		logger.Error("Redis URL not provided")
		os.Exit(1)
	}

	redisClient := internal.NewRedisClient(redisUrl)

	app := &interfaces.App{
		Logger:    *logger,
		Redis:     redisClient,
		Db:        database.GetDbInstance(),
		Koa:       koa,
		Fs:        fs,
		Constants: initConstants(),
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Start HTTP server in a goroutine
	go func() {
		defer wg.Done()
		api.InitHTTPServer(app)
	}()

	go func() {
		defer wg.Done()
		websocket_server.InitWebsocketServer(app)
	}()
	wg.Wait()
	logger.Info("Application ready!!")

}
