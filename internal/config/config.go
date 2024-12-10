package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env              string        `env:"ENV" env-required:"true"`
	DbPath           string        `env:"DB_PATH" env-required:"true"`
	HttpAddr         string        `env:"HTTP_ADDR" env-required:"true"`
	HttpReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"10s"`
	HttpWriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("unable to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
