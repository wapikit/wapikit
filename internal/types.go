package internal

import (
	"database/sql"
	"log/slog"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/stuffbin"
	"github.com/labstack/echo/v4"
)

type PermissionRole string

type constants struct {
	SiteName      string `koanf:"site_name"`
	RootURL       string `koanf:"root_url"`
	LogoURL       string `koanf:"logo_url"`
	FaviconURL    string `koanf:"favicon_url"`
	AdminUsername []byte `koanf:"admin_username"`
	AdminPassword []byte `koanf:"admin_password"`
	IsDevelopment bool
	IsProduction  bool
}

type App struct {
	Db        *sql.DB
	Logger    slog.Logger
	Koa       *koanf.Koanf
	Fs        stuffbin.FileSystem
	Constants *constants
}

const (
	SuperAdmin PermissionRole = "superadmin"
	AdminRole  PermissionRole = "admin"
	UserRole   PermissionRole = "user"
)

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
	echo.Context `json:",inline"`
	App          App            `json:"app,omitempty"`
	Session      ContextSession `json:"session,omitempty"`
}
