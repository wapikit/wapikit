package main

import (
	"github.com/labstack/echo/v4"
)

// registerHandlers registers HTTP handlers.
func initHTTPHandlers(e *echo.Echo, app *App) {
	// Group of private handlers with BasicAuth.
	g := e.Group("/")

	// ! TODO: add the auth JWT middleware
	// g.Use()

	g.GET("/", func(ctx echo.Context) error {
		return ctx.String(200, "OK")
	})

}
