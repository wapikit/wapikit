package interfaces

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Route struct {
	Path                    string         `json:"path"`
	Method                  string         `json:"method"`
	PermissionRoleLevel     PermissionRole `json:"permissionRoleLevel"` // say level is superAdmin so only super admin can access this route, but if level is user role then all the roles above the user role which is super admin and admins can access this route
	Handler                 func(context CustomContext) error
	IsAuthorizationRequired bool
}

type PermissionRole string

type ApiService interface {
	Register(server *echo.Echo)
}

type Handler interface {
	Handle(context echo.Context) error
}

type EchoHandler func(context echo.Context) error

func (eh EchoHandler) Handle(context echo.Context) error {
	return eh(context)
}

type CustomHandler func(context CustomContext) error

func (ch CustomHandler) Handle(context echo.Context) error {
	session := context.Get("Session").(ContextSession)
	app := context.Get("App").(*App)
	if session != (ContextSession{}) {
		return ch(
			CustomContext{
				Context: context,
				App:     *app,
				Session: session,
			},
		)
	} else {
		return ch(
			CustomContext{
				Context: context,
				App:     *app,
				Session: ContextSession{},
			},
		)
	}

}

const (
	SuperAdmin PermissionRole = "superadmin"
	AdminRole  PermissionRole = "admin"
	UserRole   PermissionRole = "user"
)

type ContextUser struct {
	UniqueId       string         `json:"unique_id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	Role           PermissionRole `json:"role"`
	OrganizationId string         `json:"organization_id"`
}

type ContextSession struct {
	Token string      `json:"token"`
	User  ContextUser `json:"user"`
}

type CustomContext struct {
	echo.Context `json:",inline"`
	App          App            `json:"app,omitempty"`
	Session      ContextSession `json:"session,omitempty"`
}

type JwtPayload struct {
	ContextUser        `json:",inline"`
	jwt.StandardClaims `json:",inline"`
}
