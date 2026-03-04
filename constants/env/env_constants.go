package env_constants

// define all the env constants here

type EnvMode string

const (
	EnvModeDevelopment EnvMode = "development"
	EnvModeStaging     EnvMode = "staging"
	EnvModeProduction  EnvMode = "production"
)
