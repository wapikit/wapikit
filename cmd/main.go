package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
	wapi "github.com/wapikit/wapi.go/pkg/client"
	api "github.com/wapikit/wapikit/api/cmd"
	cache "github.com/wapikit/wapikit/internal/core/redis"
	"github.com/wapikit/wapikit/internal/database"
	"github.com/wapikit/wapikit/internal/interfaces"
	campaign_manager "github.com/wapikit/wapikit/manager/campaign"
	websocket_server "github.com/wapikit/wapikit/websocket-server"
)

// because this will be a single binary, we will be providing the flags here
// 1. --install to install the setup the app, but it will be idempotent
// 2. --migrate to apply the migration to the database
// 3. --config to setup the config files
// 4. --version to check the version of the application running
// 5. --help to check the

var (
	// Global variables
	logger             = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	koa                = koanf.New(".")
	fs                 stuffbin.FileSystem
	appDir             string = "."
	frontendDir        string = "frontend/out"
	isDebugModeEnabled bool
)

func init() {
	initFlags()

	if koa.Bool("version") {
		logger.Info("current version of the application")
	}

	if koa.Bool("debug") {
		isDebugModeEnabled = true
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	// Generate new config.
	if koa.Bool("new-config") {
		path := koa.Strings("config")[0]
		if err := newConfigFile(path); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		logger.Debug("generated %s. Edit and run --install", path, nil)
		os.Exit(0)
	}

	// here appDir is for config file packing, frontendDir is for the frontend built output and static dir is any other static files and the public
	fs = initFS(appDir, frontendDir)
	loadConfigFiles(koa.Strings("config"), koa)

	if koa.Bool("install") {
		logger.Info("Installing the application")
		// ! should be idempotent
		installApp(database.GetDbInstance(koa.String("database.url")), fs, !koa.Bool("yes"), koa.Bool("idempotent"))
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

	redisUrl := koa.String("redis.url")

	if redisUrl == "" {
		logger.Error("Redis URL not provided")
		os.Exit(1)
	}

	fmt.Println("Redis URL: ", redisUrl)

	redisClient := cache.NewRedisClient(redisUrl)
	dbInstance := database.GetDbInstance(koa.String("database.url"))

	wapiClient := wapi.New(&wapi.ClientConfig{
		BusinessAccountId: koa.String("whatsapp.businessAccountId"),
		ApiAccessToken:    koa.String("whatsapp.apiKey"),
		WebhookSecret:     koa.String("whatsapp.webhookSecret"),
	})

	app := &interfaces.App{
		Logger:          *logger,
		Redis:           redisClient,
		Db:              dbInstance,
		Koa:             koa,
		Fs:              fs,
		Constants:       initConstants(),
		CampaignManager: campaign_manager.NewCampaignManager(dbInstance, *logger),
		WapiClient:      wapiClient,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// * indefinitely run the campaign manager
	go app.CampaignManager.Run()

	// Start HTTP server in a goroutine
	go func() {
		defer wg.Done()
		api.InitHTTPServer(app)
	}()

	go func() {
		defer wg.Done()
		websocket_server.InitWebsocketServer(app, &wg)
	}()

	wg.Wait()
	logger.Info("Application ready!!")

}
