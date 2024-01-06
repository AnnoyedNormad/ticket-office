package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	LogLevel    string `yaml:"log_level" env-default:"dev"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Addr         string        `yaml:"addr" env-default:"localhost:8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"10s"`
}

func MustLoad() *Config {
	op := "config.MustLoad"

	configPath := os.Getenv("CONFIG_PATH")

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("%v: %v", op, err)
	}

	return &cfg
}
