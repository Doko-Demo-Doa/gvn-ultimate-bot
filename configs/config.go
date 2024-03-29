package configs

import "os"

const (
	prod = "production"
)

// Config object
type Config struct {
	Env       string         `env:"DEPLOYMENT_ENV"`
	Pepper    string         `env:"PEPPER"`
	HMACKey   string         `env:"HMAC_KEY"`
	Postgres  PostgresConfig `json:"postgres"`
	JWTSecret string         `env:"JWT_SIGN_KEY"`
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
		Pepper:    os.Getenv("PEPPER"),
		HMACKey:   os.Getenv("HMAC_KEY"),
		Postgres:  GetPostgresConfig(),
		JWTSecret: os.Getenv("JWT_SIGN_KEY"),
		Host:      os.Getenv("APP_HOST"),
		Port:      os.Getenv("PORT"),
		FromEmail: os.Getenv("EMAIL_FROM"),
	}
}
