package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sarthakjdev/wapikit/api/services/analytics_service"
	"github.com/sarthakjdev/wapikit/api/services/auth_service"
	"github.com/sarthakjdev/wapikit/api/services/campaign_service"
	"github.com/sarthakjdev/wapikit/api/services/contact_service"
	"github.com/sarthakjdev/wapikit/api/services/conversation_service"
	contact_list_service "github.com/sarthakjdev/wapikit/api/services/list_service"
	"github.com/sarthakjdev/wapikit/api/services/next_files_service"
	organization_service "github.com/sarthakjdev/wapikit/api/services/orgnisation_service"
	"github.com/sarthakjdev/wapikit/api/services/system_service"
	webhook_service "github.com/sarthakjdev/wapikit/api/services/whatsapp_webhook_service"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

// initHTTPServer sets up and runs the app's main HTTP server and blocks forever.
func InitHTTPServer(app *interfaces.App) *echo.Echo {
	logger := app.Logger
	koa := app.Koa
	logger.Info("initializing HTTP server")
	var server = echo.New()
	server.HideBanner = true
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	isFrontendHostedSeparately := app.Koa.Bool("is_frontend_separately_hosted")

	logger.Info("isFrontendHostedSeparately: %v", isFrontendHostedSeparately)

	if !isFrontendHostedSeparately {
		// we want to mount the next.js output to "/" , i.e, / -> "index.html" , /about -> "about.html"
		fileServer := app.Fs.FileServer()
		server.GET("/*", echo.WrapHandler(fileServer))
	}

	// Mounting all HTTP handlers.
	mountHandlerServices(server, app)

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

// registerHandlers registers HTTP handlers.
func mountHandlerServices(e *echo.Echo, app *interfaces.App) {
	logger := app.Logger
	constants := app.Constants
	koa := app.Koa
	isFrontendHostedSeparately := koa.Bool("is_frontend_separately_hosted")
	corsOrigins := []string{}

	if constants.IsDevelopment {
		corsOrigins = append(corsOrigins, "http://localhost:3000")
	} else if constants.IsProduction {
		corsOrigins = append(corsOrigins, koa.String("app.cors_allowed_origins"))
	} else {
		panic("invalid environment")
	}

	app.Logger.Info("corsOrigins: %v", corsOrigins)

	// logger middleware
	e.Use(middleware.Logger())

	// compression middleware
	e.Use(middleware.Gzip())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     corsOrigins,
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderCacheControl, "x-access-token"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	servicesToRegister := []interfaces.ApiService{}
	authService := auth_service.NewAuthService()
	organizationService := organization_service.NewOrganizationService()
	campaignService := campaign_service.NewCampaignService()
	analyticsService := analytics_service.NewAnalyticsService()
	contactsService := contact_service.NewContactService()
	conversationService := conversation_service.NewConversationService()
	contactListService := contact_list_service.NewContactListService()
	systemService := system_service.NewSystemService()
	whatsappWebhookService := webhook_service.NewWhatsappWebhookServiceWebhookService()

	servicesToRegister = append(
		servicesToRegister,
		authService,
		campaignService,
		contactListService,
		contactsService,
		conversationService,
		systemService,
		whatsappWebhookService,
		analyticsService,
		organizationService,
	)

	if !isFrontendHostedSeparately {
		logger.Info("Frontend is not hosted separately")
		nextFileServerService := next_files_service.NewNextFileServerService()
		servicesToRegister = append(servicesToRegister, nextFileServerService)
	}

	// ! TODO: check for feature flags here

	for _, service := range servicesToRegister {
		service.Register(e)
	}
}
