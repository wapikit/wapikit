package main

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
	api "github.com/wapikit/wapikit/api/cmd"

	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/internal/database"
	"github.com/wapikit/wapikit/internal/services/ai_service"
	"github.com/wapikit/wapikit/internal/services/encryption_service"
	cache_service "github.com/wapikit/wapikit/internal/services/redis_service"

	campaign_manager "github.com/wapikit/wapikit/internal/campaign_manager"
	websocket_server "github.com/wapikit/wapikit/websocket-server"
)

// because this will be a single binary, we will be providing the flags here
// 1. --install to install the setup the app, but it will be idempotent
// 3. --config to setup the config files
// 4. --version to check the version of the application running
// 5. --help to check the
// 6. --debug to enable the debug mode
// 7. --new-config to generate a new config file
// 8. --yes to assume 'yes' to prompts during --install/upgrade
// 9. --ws to start the websocket server
// 10. --server to start the API server // can run multiple instance, is stateless
// 11. --cm to start the campaign manager // should run only one instance at any point of time

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

	// here appDir is for config file packing, frontendDir is for the frontend built output

	// ! TODO: find a fix because this is not going to work in the single binary mode
	fs = initFS(appDir, "")
	loadConfigFiles(koa.Strings("config"), koa)

	// load environment variables, configs can also be loaded using the environment variables, using prefix WAPIKIT_
	// for example, WAPIKIT_redis__url is equivalent of redis.url as in config.toml
	if err := koa.Load(env.Provider("WAPIKIT_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "WAPIKIT_")), "__", ".", -1)
	}), nil); err != nil {
		logger.Error("error loading config from env: %v", err, nil)
	}

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

	redisClient := cache_service.NewRedisClient(redisUrl)
	dbInstance := database.GetDbInstance(koa.String("database.url"))

	constants := initConstants()

	app := &interfaces.App{
		Logger:          *logger,
		Redis:           redisClient,
		Db:              dbInstance,
		Koa:             koa,
		Fs:              fs,
		Constants:       constants,
		CampaignManager: campaign_manager.NewCampaignManager(dbInstance, *logger),
	}

	app.EncryptionService = encryption_service.NewEncryptionService(
		logger,
		koa.String("app.encryption_key"),
	)

	if constants.IsCloudEdition {
		//  ! if this is cloud edition, then we can initiate certain service here only
		aiService := ai_service.NewAiService(
			logger,
			redisClient,
			dbInstance,
			koa.String("ai.api_key"),
		)
		app.AiService = aiService
	}

	var wg sync.WaitGroup
	wg.Add(3)

	doStartWebsocketServer := koa.Bool("ws")
	doStartAPIServer := koa.Bool("server")
	doStartCampaignManager := koa.Bool("cm")

	isSingleBinaryMode := !doStartWebsocketServer && !doStartAPIServer && !doStartCampaignManager

	if isSingleBinaryMode {
		doStartAPIServer = true
		doStartWebsocketServer = true
		doStartCampaignManager = true
	}

	if doStartCampaignManager {
		// * indefinitely run the campaign manager
		go app.CampaignManager.Run()
	}

	if doStartAPIServer {
		go func() {
			defer wg.Done()
			api.InitHTTPServer(app)
		}()
	}

	if doStartWebsocketServer {
		go func() {
			defer wg.Done()
			websocket_server.InitWebsocketServer(app, &wg)
		}()
	}

	wg.Wait()
}
