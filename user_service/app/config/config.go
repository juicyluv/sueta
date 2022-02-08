package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
)

// Config describes an application configuration structure.
type Config struct {
	Http struct {
		Port           string `yaml:"port" env-default:"8080"`
		MaxHeaderBytes int    `yaml:"maxHeaderBytes" env-default:"1"`
		ReadTimeout    int    `yaml:"readTimeout" env-default:"20"`
		WriteTimeout   int    `yaml:"writeTimeout" env-default:"20"`
		RequestTimeout int    `yaml:"requestTimeout" env-default:"15"`
	} `yaml:"http" env-required:"true"`
	DB struct {
		URL      string `env:"MONGO_URL" env-required:"true"`
		Database string `yaml:"database" env-required:"true"`
	} `yaml:"mongo" env-required:"true"`
	RedisDSN string `env:"REDIS_DSN" env-required:"true"`
}

var instance *Config
var once sync.Once

// Get loads .env file and config from given path.
// Returns config instance if everything is OK
// or an error if something went wrong.
func Get(configPath string) *Config {
	logger := logger.GetLogger()

	logger.Info("loading .env file")
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("could not load .env file: %v", err)
	}
	logger.Info("loaded .env file")

	once.Do(func() {
		logger.Info("reading application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	logger.Info("done reading application config")

	return instance
}
