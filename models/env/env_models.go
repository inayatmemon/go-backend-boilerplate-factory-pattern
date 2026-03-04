package env_models

import (
	env_constants "go_boilerplate_project/constants/env"
)

type Environment struct {
	App       *AppEnv
	Databases *DatabasesEnv
	Logger    *LoggerEnv
}

type AppEnv struct {
	EnvMode    *env_constants.EnvMode
	ApiBaseUrl string
	AppName    string
	AppVersion string
	AppPort    int
	AppHost    string
}

type DatabasesEnv struct {
	MongoDB *MongoDB
	Redis   *Redis
	MySQL   *MySQL
}

type MongoDB struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type Redis struct {
	Host     string
	Port     int
	Database int
	Username string
	Password string
}

type MySQL struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type LoggerEnv struct {
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}
