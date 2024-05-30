package main

import (
	"fmt"
	"net/http"

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
			return echo.ErrForbidden
		}
	}
}

// registerHandlers registers HTTP handlers.
func mountHandlers(e *echo.Echo, app *App) {
	// ! TODO: enable cors
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"", ""},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// }))

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
