package plugin

import "github.com/labstack/echo/v4"

type PluginBaseConfig struct {
	Name        string `json:"name"`
	RestApiPath string `json:"rest_api_path"`
}

type BasePluginInterface interface {
	Register(server echo.Echo) error
	Bootstrap(server echo.Echo) error
	GetConfig() PluginBaseConfig
}

type BasePlugin struct {
	Name string `json:"name"`
}

func (p *BasePlugin) Register() error {
	return nil
}
