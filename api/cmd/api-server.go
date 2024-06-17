package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/knadh/stuffbin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sarthakjdev/wapikit/api/services/auth_service"
	"github.com/sarthakjdev/wapikit/api/services/campaign_service"
	"github.com/sarthakjdev/wapikit/api/services/contact_service"
	"github.com/sarthakjdev/wapikit/api/services/conversation_service"
	contact_list_service "github.com/sarthakjdev/wapikit/api/services/list_service"
	"github.com/sarthakjdev/wapikit/api/services/next_files_service"
	"github.com/sarthakjdev/wapikit/api/services/organization_member_service"
	"github.com/sarthakjdev/wapikit/api/services/system_service"
	webhook_service "github.com/sarthakjdev/wapikit/api/services/whatsapp_webhook_service"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

// initHTTPServer sets up and runs the app's main HTTP server and blocks forever.
func initHTTPServer(app *App) *echo.Echo {
	app.logger.Info("initializing HTTP server")
	var server = echo.New()
	logger := app.logger
	server.HideBanner = true
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	isFrontendHostedSeparately := app.koa.Bool("is_frontend_separately_hosted")

	if !isFrontendHostedSeparately {
		// we want to mount the next.js output to "/" , i.e, / -> "index.html" , /about -> "about.html"
		fileServer := app.fs.FileServer()
		tpl, err := stuffbin.ParseTemplatesGlob(initTplFuncs(app.constants), app.fs, "/*.html")
		if err != nil {
			logger.Error("error parsing public templates: %v", err)
		}

		server.Renderer = &tplRenderer{
			templates:  tpl,
			SiteName:   app.constants.SiteName,
			RootURL:    app.constants.RootURL,
			LogoURL:    app.constants.LogoURL,
			FaviconURL: app.constants.FaviconURL,
		}

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
func mountHandlerServices(e *echo.Echo, app *App) {
	isFrontendHostedSeparately := app.koa.Bool("is_frontend_separately_hosted")
	corsOrigins := []string{}

	if app.constants.IsDevelopment {
		corsOrigins = append(corsOrigins, koa.String("address"))
	} else if app.constants.IsProduction {
		corsOrigins = append(corsOrigins, koa.String("cors_allowed_origins"))
	} else {
		panic("invalid environment")
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	servicesToRegister := []interfaces.ApiService{}

	organizationMemberService := organization_member_service.NewOrganizationMemberService()
	authService := auth_service.NewAuthService()
	campaignService := campaign_service.NewCampaignService()
	// analyticsService := analytics_service.NewAnalyticsService()
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
		organizationMemberService,
		systemService,
		whatsappWebhookService,
	)

	if !isFrontendHostedSeparately {
		nextFileServerService := next_files_service.NewNextFileServerService()
		servicesToRegister = append(servicesToRegister, nextFileServerService)
	}

	// ! TODO: check for feature flags here

	for _, service := range servicesToRegister {
		service.Register(e)
	}
}
