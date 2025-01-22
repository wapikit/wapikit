package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wapikit/wapikit/interfaces"
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

	isFrontendHostedSeparately := app.Koa.Bool("app.is_frontend_separately_hosted")

	if !isFrontendHostedSeparately && app.Constants.IsProduction {
		// we want to mount the next.js output to "/" , i.e, / -> "index.html" , /about -> "about.html"
		fileServer := app.Fs.FileServer()
		server.GET("/*", echo.WrapHandler(fileServer))
	}

	addMiddlewares(server, app)

	// Mounting all HTTP handlers.
	mountHandlerServices(server, app)

	// getting th server address from config and falling back to localhost:8000
	serverAddress := koa.String("app.address")

	if serverAddress == "" {
		serverAddress = "localhost:8000"
	}

	// Start the server.
	func() {
		logger.Info("starting HTTP server on %s", serverAddress, nil) // Add a placeholder value as the final argument
		if err := server.Start(serverAddress); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("HTTP server shut down")
			} else {
				logger.Error("error starting HTTP server: %v", err.Error(), nil)
			}
		}
	}()

	return server
}

func addMiddlewares(e *echo.Echo, app *interfaces.App) {
	constants := app.Constants
	logger := app.Logger
	koa := app.Koa

	var origins []string

	corsAllowedOrigins := koa.String("app.cors_allowed_origins")
	if err := json.Unmarshal([]byte(corsAllowedOrigins), &origins); err != nil {
		// If unmarshalling fails, try to parse it as a TOML array
		if strings.HasPrefix(corsAllowedOrigins, "[") && strings.HasSuffix(corsAllowedOrigins, "]") {
			corsAllowedOrigins = strings.TrimPrefix(corsAllowedOrigins, "[")
			corsAllowedOrigins = strings.TrimSuffix(corsAllowedOrigins, "]")
			origins = strings.Split(corsAllowedOrigins, " ")
			for i, origin := range origins {
				logger.Debug("allowing origin", origin, nil)
				origins[i] = strings.TrimSpace(strings.Trim(origins[i], `"`))
			}
		} else {
			fmt.Println("Error parsing CORS allowed origins:", err)
			return
		}
	}

	// logger middleware
	if constants.IsDebugModeEnabled {
		e.Use(middleware.Logger())
	}

	// compression middleware
	e.Use(middleware.Gzip())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     origins,
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderCacheControl, "x-access-token"},
		AllowMethods:     []string{http.MethodPost, http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions},
		MaxAge:           5,
	}))

}
