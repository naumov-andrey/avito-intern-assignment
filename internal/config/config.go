package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

type PostgresConfig struct {
	User     string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
	Password string
}

type HTTPConfig struct {
	Port string
}

type Config struct {
	Postgres PostgresConfig
	HTTP     HTTPConfig
}

func Init(cfgPath string, cfgName string) (*Config, error) {
	viper.AddConfigPath(cfgPath)
	viper.SetConfigName(cfgName)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return nil, err
	}
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	cfg.Postgres.Password = os.Getenv("DB_PASSWORD")

	return &cfg, nil
}
