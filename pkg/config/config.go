package config

import (
	"github.com/danryan/env"
	"github.com/gplume/todo-list/pkg/utils"
)

// Config structure for app options
type Config struct {
	FrameWorkMode string `env:"key=FRAMEWORK_MODE default=debug"`
	UsageMode     string `env:"key=USAGE_MODE default=dev"`
	BaseDir       string `env:"key=BASE_DIR"`
	AppName       string `env:"key=APP_NAME"`
	AppDomain     string `env:"key=APP_DOMAIN required=true"`
	ServerPort    int    `env:"key=PORT default=8000"`
	DBType        string `env:"key=DB_TYPE required=true"`
	DBDirectory   string `env:"key=DB_DIRECTORY default=db"`
	DBName        string `env:"key=DB_NAME required=true"`
	DBTestName    string `env:"key=DB_TEST_NAME required=true"`
}

// New return a new config structure
func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Process(cfg); err != nil {
		return cfg, err
	}

	// Absolute App Dir
	if cfg.BaseDir == "" {
		cfg.BaseDir = utils.GetDefaultBaseDir()
	}

	// do others customizing of *config here

	return cfg, nil
}
