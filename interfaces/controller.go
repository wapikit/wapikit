package interfaces

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/internal/api_types"
)

type RateLimitConfig struct {
	MaxRequests    int   `json:"maxRequests"`
	WindowTimeInMs int64 `json:"windowTime"`
}

type RouteMetaData struct {
	PermissionRoleLevel  api_types.UserPermissionLevelEnum `json:"permissionRoleLevel"`
	RequiredPermission   []api_types.RolePermissionEnum    `json:"requiredPermission"`
	RateLimitConfig      RateLimitConfig                   `json:"rateLimitConfig"`
	RequiredFeatureFlags []string                          `json:"requiredFeatureFlags"`
}

type Route struct {
	Path                    string `json:"path"`
	Method                  string `json:"method"`
	Handler                 Handler
	IsAuthorizationRequired bool
	MetaData                RouteMetaData `json:"metaData"`
}

type ApiController interface {
	Register(server *echo.Echo)
	GetControllerName() string
}

type Handler interface {
	Handle(context echo.Context) error
}

type HandlerWithoutSession func(context ContextWithoutSession) error

func (eh HandlerWithoutSession) Handle(context echo.Context) error {
	return eh(context.(ContextWithoutSession))
}

type HandlerWithSession func(context ContextWithSession) error

func (ch HandlerWithSession) Handle(context echo.Context) error {
	return ch(context.(ContextWithSession))
}

type ContextUser struct {
	Name           string                            `json:"name"`
	UniqueId       string                            `json:"unique_id"`
	Username       string                            `json:"username"`
	Email          string                            `json:"email"`
	Role           api_types.UserPermissionLevelEnum `json:"role"`
	OrganizationId string                            `json:"organization_id"`
}

type ContextSession struct {
	Token string      `json:"token"`
	User  ContextUser `json:"user"`
}

type ContextWithSession struct {
	echo.Context `json:",inline"`
	App          App            `json:"app,omitempty"`
	Session      ContextSession `json:"session,omitempty"`
}

type ContextWithoutSession struct {
	echo.Context `json:",inline"`
	App          App `json:"app,omitempty"`
}

type JwtPayload struct {
	ContextUser        `json:",inline"`
	jwt.StandardClaims `json:",inline"`
}
