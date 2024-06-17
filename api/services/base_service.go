package services

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type BaseService struct {
	Name        string `json:"name"`
	RestApiPath string `json:"rest_api_path"`
	Routes      []interfaces.Route
}

func (s *BaseService) GetServiceName() string {
	return s.Name

}

func (s *BaseService) GetRoutes() []interfaces.Route {
	return s.Routes
}

func (s *BaseService) GetRestApiPath() string {
	return s.RestApiPath
}

func isAuthorized(role interfaces.PermissionRole, routerPermissionLevel interfaces.PermissionRole) bool {
	switch role {
	case interfaces.SuperAdmin:
		return true
	case interfaces.AdminRole:
		return routerPermissionLevel == interfaces.AdminRole || routerPermissionLevel == interfaces.UserRole
	case interfaces.UserRole:
		return routerPermissionLevel == interfaces.UserRole
	default:
		return false
	}
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {

		app := context.Get("app").(*interfaces.App)
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
		metadata := context.Get("metadata").(interfaces.Route)

		// ! TODO: fetch the user from db and check role here
		mockRole := interfaces.SuperAdmin

		if isAuthorized(mockRole, metadata.PermissionRoleLevel) {
			return next(interfaces.CustomContext{
				Context: context,
				App:     *app,
				Session: interfaces.ContextSession{
					Token: authToken,
					User: interfaces.ContextUser{
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

func rateLimiter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		// rate limit the request
		// return error if rate limit is exceeded
		return next(context)
	}
}

// Register function now uses the Routes field
func (service *BaseService) Register(server *echo.Echo) {
	group := server.Group(service.RestApiPath)
	for _, route := range service.Routes {
		handler := interfaces.CustomHandler(route.Handler).Handle
		if route.IsAuthorizationRequired {
			handler = authMiddleware(interfaces.CustomHandler(route.Handler).Handle)
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
