package main

import "github.com/danryan/env"

type config struct {
	FrameWorkMode string `env:"key=FRAMEWORK_MODE default=debug"`
	UsageMode     string `env:"key=USAGE_MODE default=dev"`
	BaseDir       string `env:"key=BASE_DIR"`
	AppName       string `env:"key=APP_NAME"`
	AppDomain     string `env:"key=APP_DOMAIN required=true"`
	ServerPort    int    `env:"key=PORT default=8000"`
	SSLEnabled    bool   `env:"key=SSL_ENABLED default=false"`
	SSLPub        string `env:"key=SSL_PUB"`
	SSLKey        string `env:"key=SSL_KEY"`
	DBType        string `env:"key=DB_TYPE required=true"`
	DBDirectory   string `env:"key=DB_DIRECTORY default=db"`
	DBName        string `env:"key=DB_NAME required=true"`
}

func newConfig() (*config, error) {
	cfg := &config{}
	if err := env.Process(cfg); err != nil {
		return cfg, err
	}

	// Absolute App Dir
	if cfg.BaseDir == "" {
		cfg.BaseDir = getDefaultBaseDir()
	}

	return cfg, nil
}
