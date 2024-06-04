package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/handlers"
)

type PermissionRole string

const (
	SuperAdmin PermissionRole = "superadmin"
	AdminRole  PermissionRole = "admin"
	UserRole   PermissionRole = "user"
)

// Helper functions to attach metadata and middleware
func CreateHandler(metadata Route, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		context.Set("metadata", metadata) // Store metadata in context
		return handler(context)
	}
}

type Route struct {
	Path                    string         `json:"path"`
	Method                  string         `json:"method"`
	PermissionRoleLevel     PermissionRole `json:"permissionRoleLevel"` // say level is superAdmin so only super admin can access this route, but if level is user role then all the roles above the user role which is super admin and admins can access this route
	Handler                 func(c echo.Context) error
	IsAuthorizationRequired bool
}

type ContextUser struct {
	UniqueId string         `json:"unique_id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Role     PermissionRole `json:"role"`
}

type ContextSession struct {
	Token string      `json:"token"`
	User  ContextUser `json:"user"`
}

type CustomContext struct {
	echo.Context
	Session ContextSession `json:"session,omitempty"`
}

func isAuthorized(role PermissionRole, routerPermissionLevel PermissionRole) bool {
	switch role {
	case SuperAdmin:
		return true
	case AdminRole:
		return routerPermissionLevel == AdminRole || routerPermissionLevel == UserRole
	case UserRole:
		return routerPermissionLevel == UserRole
	default:
		return false
	}
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		headers := context.Request().Header
		authToken := headers.Get("x-access-token")
		origin := headers.Get("origin")

		if authToken == "" {
			return echo.ErrUnauthorized
		}

		// verify the token
		// get the user
		// inject user to context
		// consider user permission level too
		// return error if token is invalid
		fmt.Println("authMiddleware: ", authToken, origin)
		metadata := context.Get("metadata").(Route)

		// ! TODO: fetch the user from db and check role here
		mockRole := SuperAdmin

		if isAuthorized(mockRole, metadata.PermissionRoleLevel) {
			return next(CustomContext{
				context,
				ContextSession{
					Token: authToken,
					User: ContextUser{
						"",
						"",
						"",
						mockRole,
					},
				},
			})
		} else {
			fmt.Println("authMiddleware: ", authToken, origin, metadata.PermissionRoleLevel, mockRole)
			return echo.ErrForbidden
		}
	}
}

func ServerMediaFiles(c echo.Context) error {
	app := c.Get("app").(*App)
	routePath := c.Request().URL.Path
	b, err := app.fs.Read(routePath)
	if err != nil {
		if err.Error() == "file does not exist" {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.Blob(http.StatusOK, mimetype.Detect(b).String(), b)
}

func serverHtmlAndNonJsAndCssFiles(c echo.Context) error {
	app := c.Get("app").(*App)
	routePath := c.Request().URL.Path
	fmt.Println("routePath: ", routePath, path.Ext(routePath))
	// check if the request is for some extension other than html or no extension
	requestedFileExt := path.Ext(routePath)
	if routePath != "/" && requestedFileExt != "" && requestedFileExt != ".html" {
		logger.Info("serving static files: %v", routePath, nil)
		b, err := app.fs.Read(routePath)
		if err != nil {
			logger.Error("error reading static file: %v", err)
			if err.Error() == "file does not exist" {
				_404File, err := app.fs.Read(path.Join("", "/404.html"))
				if err != nil {
					return echo.NewHTTPError(http.StatusNotFound)
				}

				return c.HTMLBlob(http.StatusOK, _404File)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.Blob(http.StatusOK, mimetype.Detect(b).String(), b)
	}

	if routePath == "/" {
		logger.Info("serving index.html")
		routePath = "/index"
	}

	b, err := app.fs.Read(path.Join("", routePath+".html"))
	if err != nil {
		logger.Error("error reading static file in end block: %v", err)

		if err.Error() == "file does not exist" {
			_404File, err := app.fs.Read(path.Join("", "/404.html"))
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return c.HTMLBlob(http.StatusOK, _404File)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.HTMLBlob(http.StatusOK, b)
}

func handleNextStaticJsAndCssRoute(c echo.Context) error {
	app := c.Get("app").(*App)
	b, err := app.fs.Read(c.Request().URL.Path)

	if err != nil {
		if err.Error() == "file does not exist" {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// check if the file is a js file or css file
	if path.Ext(c.Request().URL.Path) == ".js" {
		return c.Blob(http.StatusOK, "application/javascript", b)
	} else {
		return c.Blob(http.StatusOK, "text/css", b)
	}
}

// registerHandlers registers HTTP handlers.
func mountHandlers(e *echo.Echo, app *App) {
	// ! TODO: enable cors here on the basis frontend url is diff 
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"", ""},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// }))

	// ! TODO: check here if the frontend hosting is enabled or not, but media file directory would still be enabled
	e.GET("/media/*", ServerMediaFiles)
	e.GET("/_next/*", handleNextStaticJsAndCssRoute)
	e.GET("/*", serverHtmlAndNonJsAndCssFiles)

	// * all the backend api path should start from "/api" as "/" is for frontend files and static files
	routes := []Route{
		{Path: "/health-check", Method: "GET", Handler: handlers.GetUsers, IsAuthorizationRequired: false},
		{Path: "/users", Method: "GET", Handler: handlers.GetUsers, IsAuthorizationRequired: true, PermissionRoleLevel: SuperAdmin},
	}

	group := e.Group("/api")
	for _, route := range routes {
		handler := route.Handler
		if route.IsAuthorizationRequired {
			handler = authMiddleware(handler)
		}
		switch route.Method {
		case http.MethodGet:
			group.GET(route.Path, handler)
		case http.MethodPost:
			group.POST(route.Path, handler)
		case http.MethodPut:
			group.PUT(route.Path, handler)
		case http.MethodDelete:
			group.DELETE(route.Path, handler)
		}
	}
}
