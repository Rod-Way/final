package config

import (
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"development"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer
	Agent
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"5500"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Agent struct {
	GoroutinesNum string `yaml:"goroutinesnum" env-default:"5"`
	AgentsNum     string `yaml:"database" env-default:"2"`
}

var conf *Config
var once sync.Once

func MustLoad(logger slog.Logger) *Config {
	once.Do(func() {
		logger.Info("read application configuration")
		conf = &Config{}
		if err := cleanenv.ReadConfig(".\\config\\config.yaml", conf); err != nil {
			help, _ := cleanenv.GetDescription(conf, nil)
			logger.Info(help)
			log.Fatalf("Configuration error: %s", err)
		}
	})
	return conf
}
