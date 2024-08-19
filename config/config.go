package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Db     `yaml:"postgres"`
	Secret `yaml:"secret"`
}

type Secret struct {
	Key string `yaml:"key"`
}

type Db struct {
	Host         string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"password"`
	Db           string `yaml:"postgres-db" env:"POSTGRES_DB" env-default:"postgres"`
	DbTimeoutSec int    `yaml:"db-timeout-sec"`
}

func ReadConfig() (*Config, error) {
    cfg := Config{}
	dir, err := os.Getwd()
    if err != nil {
        return nil, fmt.Errorf("can't get current directory: %v", err)
    }

    configPath := filepath.Join(dir, "config", "config.yml")

    err = cleanenv.ReadConfig(configPath, &cfg)
    if err != nil {
        return nil, fmt.Errorf("read config error: %v", err)
    }
    err = cleanenv.ReadEnv(&cfg)
    if err != nil {
        return nil, fmt.Errorf("read env error: %v", err)
    }

    return &cfg, nil
}
