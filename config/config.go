package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBName              string `envconfig:"DBNAME" required:"true"`
	DBHost              string `envconfig:"HOST" required:"true"`
	DBuser              string `envconfig:"USER" required:"true"`
	DBPassword          string `envconfig:"PASSWORD" required:"true"`
	DBPort              int    `envconfig:"PORT" required:"true"`
	DBMaxIdle           int    `envconfig:"DB_MAX_IDLE" default:"100"`
	DBMaxConnection     int    `envconfig:"DB_MAX_CONNECTION" default:"100"`
	DBConnectionTimeout int    `envconfig:"DB_CONNECTION_TIMEOUT" default:"10"`
}

func Setup() *Config {
	var db Config
	envconfig.MustProcess("MYSQL", &db)
	return &db
}
