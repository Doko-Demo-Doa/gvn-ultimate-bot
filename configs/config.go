package configs

import "os"

const (
	prod = "production"
)

// Config object
type Config struct {
	Env       string         `env:"DEPLOYMENT_ENV"`
	Postgres  PostgresConfig `json:"postgres"`
	Host      string         `env:"APP_HOST"`
	Port      string         `env:"PORT"`
	FromEmail string         `env:"EMAIL_FROM"`
}

// IsProd Checks if env is production
func (c Config) IsProd() bool {
	return c.Env == prod
}

// GetConfig gets all config for the application
func GetConfig() Config {
	return Config{
		Env:       os.Getenv("DEPLOYMENT_ENV"),
		Postgres:  GetPostgresConfig(),
		Host:      os.Getenv("APP_HOST"),
		Port:      os.Getenv("PORT"),
		FromEmail: os.Getenv("EMAIL_FROM"),
	}
}
