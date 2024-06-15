package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
	"github.com/sarthakjdev/wapikit/services/auth_service"
	"github.com/sarthakjdev/wapikit/services/next_files_service"
)

// registerHandlers registers HTTP handlers.
func mountHandlers(e *echo.Echo, app *App) {
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

	// if frontend is hosted separately, then do not serve the frontend files
	// e.GET("/webhook", shared_interfaces.CustomHandler(wapiClient.GetWebhookGetRequestHandler()).Handle)
	// e.POST("/webhook", shared_interfaces.CustomHandler(wapiClient.GetWebhookPostRequestHandler()).Handle)

	servicesToRegister := []interfaces.ApiService{}

	authService := auth_service.NewAuthService()

	servicesToRegister = append(servicesToRegister, authService)

	if isFrontendHostedSeparately {
		nextFileServerService := next_files_service.NewNextFileServerService()
		servicesToRegister = append(servicesToRegister, nextFileServerService)
	}

	// ! TODO: check for feature flags here

	for _, service := range servicesToRegister {
		service.Register(e)
	}

}
