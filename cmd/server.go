package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sarthakjdev/wapikit/handlers"
	"github.com/sarthakjdev/wapikit/internal"
)

type Handler interface {
	Handle(context echo.Context) error
}

type EchoHandler func(context echo.Context) error

func (eh EchoHandler) Handle(context echo.Context) error {
	return eh(context)
}

type CustomHandler func(context internal.CustomContext) error

func (ch CustomHandler) Handle(context echo.Context) error {
	session := context.Get("Session").(internal.ContextSession)
	app := context.Get("App").(*internal.App)
	if session != (internal.ContextSession{}) {
		return ch(
			internal.CustomContext{
				Context: context,
				App:     *app,
				Session: session,
			},
		)
	} else {
		return ch(
			internal.CustomContext{
				Context: context,
				App:     *app,
				Session: internal.ContextSession{},
			},
		)
	}

}

type Route struct {
	Path                    string                  `json:"path"`
	Method                  string                  `json:"method"`
	PermissionRoleLevel     internal.PermissionRole `json:"permissionRoleLevel"` // say level is superAdmin so only super admin can access this route, but if level is user role then all the roles above the user role which is super admin and admins can access this route
	Handler                 func(context internal.CustomContext) error
	IsAuthorizationRequired bool
}

func isAuthorized(role internal.PermissionRole, routerPermissionLevel internal.PermissionRole) bool {
	switch role {
	case internal.SuperAdmin:
		return true
	case internal.AdminRole:
		return routerPermissionLevel == internal.AdminRole || routerPermissionLevel == internal.UserRole
	case internal.UserRole:
		return routerPermissionLevel == internal.UserRole
	default:
		return false
	}
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {

		app := context.Get("app").(*internal.App)
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
		mockRole := internal.SuperAdmin

		if isAuthorized(mockRole, metadata.PermissionRoleLevel) {
			return next(internal.CustomContext{
				Context: context,
				App:     *app,
				Session: internal.ContextSession{
					Token: authToken,
					User: internal.ContextUser{
						UniqueId: "",
						Username: "",
						Email:    "",
						Role:     mockRole,
					},
				},
			})
		} else {
			fmt.Println("authMiddleware: ", authToken, origin, metadata.PermissionRoleLevel, mockRole)
			return echo.ErrForbidden
		}
	}
}

// registerHandlers registers HTTP handlers.
func mountHandlers(e *echo.Echo, app *App) {
	logger.Info("is_frontend_separately_hosted", app.koa.Bool("is_frontend_separately_hosted"), nil)
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
	if isFrontendHostedSeparately {
		e.GET("/media/*", CustomHandler(handlers.ServerMediaFiles).Handle)
		e.GET("/_next/*", CustomHandler(handlers.HandleNextStaticJsAndCssRoute).Handle)
		e.GET("/*", CustomHandler(handlers.ServerHtmlAndNonJsAndCssFiles).Handle)
	}

	// attach the webhook handler

	wapiClient := internal.GetWapiCloudClient(
		koa.String("PHONE_NUMBER_ID"),
		koa.String("WHATSAPP_BUSINESS_ACCOUNT_ID"),
		koa.String("WHATSAPP_WEBHOOK_SECRET"),
		koa.String("WHATSAPP_API_ACCESS_TOKEN"),
	)

	e.GET("/webhook", EchoHandler(wapiClient.GetWebhookGetRequestHandler()).Handle)

	e.POST("/webhook", EchoHandler(wapiClient.GetWebhookPostRequestHandler()).Handle)

	// * all the backend api path should start from "/api" as "/" is for frontend files and static files
	routes := []Route{
		{Path: "/health-check", Method: "GET", Handler: handlers.HandleHealthCheck, IsAuthorizationRequired: false},
		{Path: "/login", Method: "POST", Handler: handlers.HandleSignIn, IsAuthorizationRequired: false, PermissionRoleLevel: internal.UserRole},
		{Path: "/members/:id", Method: "GET", Handler: handlers.GetOrgMemberById, IsAuthorizationRequired: true, PermissionRoleLevel: internal.AdminRole},
		{Path: "/members/:id", Method: "PUT", Handler: handlers.GetOrgMemberById, IsAuthorizationRequired: true, PermissionRoleLevel: internal.AdminRole},
		{Path: "/members/:id", Method: "DELETE", Handler: handlers.DeleteOrgMemberById, IsAuthorizationRequired: true, PermissionRoleLevel: internal.AdminRole},
		{Path: "/members", Method: "GET", Handler: handlers.GetAllOrganizationMembers, IsAuthorizationRequired: true, PermissionRoleLevel: internal.AdminRole},
		{Path: "/members", Method: "POST", Handler: handlers.CreateNewOrganizationMember, IsAuthorizationRequired: true, PermissionRoleLevel: internal.AdminRole},
	}

	group := e.Group("/api")
	for _, route := range routes {
		handler := CustomHandler(route.Handler).Handle
		if route.IsAuthorizationRequired {
			handler = authMiddleware(CustomHandler(route.Handler).Handle)
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
