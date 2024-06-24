package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Database
	Server
}

type Server struct {
	Address      string        `yaml:"address"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type Database struct {
	Name    string `yaml:"name"`
	ConnStr string `yaml:"conn_str"`
}

func Load() *Config {

	//Maybe get config path from env
	configPath := "config/local.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exists: %s", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return cfg
}
