package configs

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string `env:"PGHOST"`
	Port     int    `env:"PGPORT"`
	User     string `env:"PGUSER"`
	Password string `env:"PGPASSWORD"`
	Name     string `env:"PGDATABASE"`
}

// Dialects returns "postgres"
func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) GetPostgresConfigInfo() gorm.Dialector {
	strconn := ""

	if c.Password == "" {
		strconn = fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s",
			c.Host, c.Port, c.User, c.Name,
		)
		return postgres.Open(strconn)
	}

	strconn = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		c.Host, c.Port, c.User, c.Password, c.Name)

	return postgres.Open(strconn)
}

func GetPostgresConfig() PostgresConfig {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	return PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}
