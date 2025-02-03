package main

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/toml"
	file "github.com/knadh/koanf/providers/file"
	posflag "github.com/knadh/koanf/providers/posflag"

	_ "github.com/go-jet/jet/v2/postgres"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
	_ "github.com/wapikit/wapikit/.db-generated/model"
	_ "github.com/wapikit/wapikit/.db-generated/table"
	"github.com/wapikit/wapikit/interfaces"
)

func initConstants() *interfaces.Constants {
	var c interfaces.Constants

	if err := koa.Unmarshal("app", &c); err != nil {
		logger.Error("error loading app config: %v", err.Error(), nil)
	}

	if koa.String("environment") == "development" {
		c.IsDevelopment = true
		c.IsProduction = false
	} else {
		c.IsProduction = true
		c.IsDevelopment = false
	}

	c.RedisEventChannelName = "ApiServerEvents"
	c.IsDebugModeEnabled = isDebugModeEnabled
	c.IsCloudEdition = koa.Bool("is_cloud_edition")
	c.IsSingleBinaryMode = koa.Bool("is_single_binary_mode")
	c.IsCommunityEdition = !c.IsCloudEdition

	return &c
}

func initFlags() {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println("Usage: wapikit [flags]:\n", f.FlagUsages())
		os.Exit(0)
	}

	// Register the command line flags.
	f.StringSlice("config", []string{"config.toml"},
		"path to one or more config files (will be merged in order)")
	f.Bool("install", false, "setup database (first time)")
	f.Bool("debug", false, "enable debug mode")
	f.Bool("new-config", false, "generate a new config file")
	f.Bool("idempotent", false, "make --install run only if the database isn't already setup")
	f.Bool("yes", false, "assume 'yes' to prompts during --install/upgrade")
	f.Bool("server", false, "starts the API server")
	f.Bool("cm", false, "starts the campaign manager")

	// ! TODO: implement and enable the below flags
	// f.Bool("upgrade", false, "upgrade database to the current version")
	// f.Bool("version", false, "show current version of the build")
	// f.Bool("db-migrate", false, "apply database migrations")
	// f.Bool("db-apply", false, "apply database migrations")

	if err := f.Parse(os.Args[1:]); err != nil {
		logger.Error("error loading flags: %v", err)
	}

	if err := koa.Load(posflag.Provider(f, ".", koa), nil); err != nil {
		logger.Error("error loading config: %v", err)
	}
}

func loadConfigFiles(filePaths []string, koa *koanf.Koanf) {
	for _, filePath := range filePaths {
		if err := koa.Load(file.Provider(filePath), toml.Parser()); err != nil {
			if os.IsNotExist(err) {
				logger.Error("config file not found. If there isn't one yet, run --config to generate one.")
			}
			logger.Error("error loading config from file: %v.", err)
		}
		logger.Info("loaded config file.")
		logger.Debug("loaded config file: %s", "filePath", filePath)
	}
}

func newConfigFile(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists. Edit it or remove it to generate a new one.", path)
	}

	// Initialize the static file system with empty appDir and frontendDir, to load the config.toml.sample.
	fs := initFS("", "")
	b, err := fs.Read("config.toml.sample")
	if err != nil {
		return fmt.Errorf("error reading sample config (is binary stuffed?): %v", err)
	}

	return os.WriteFile(path, b, 0644)
}
