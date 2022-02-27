package configs

type Config struct {
	Env       string         `env:"ENV"`
	Pepper    string         `env:"PEPPER"`
	HMACKey   string         `env:"HMAC_KEY"`
	Postgres  PostgresConfig `json:"postgres"`
	JWTSecret string         `env:"JWT_SIGN_KEY"`
	Host      string         `env:"APP_HOST"`
	Port      string         `env:"APP_PORT"`
	FromEmail string         `env:"EMAIL_FROM"`
}
