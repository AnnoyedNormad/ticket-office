package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	LogLevel    string `yaml:"log_level" env-default:"dev"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Addr         string `yaml:"addr" env-default:"localhost:8080"`
	ReadTimeout  string `yaml:"read_timeout" env-default:"10"`
	WriteTimeout string `yaml:"write_timeout" env-default:"10"`
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
