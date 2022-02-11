package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
)

// Config describes an application configuration structure.
type Config struct {
	// Http represents configuration for http server.
	Http struct {
		Port           string `yaml:"port" env-default:"8080"`
		MaxHeaderBytes int    `yaml:"maxHeaderBytes" env-default:"1"`
		ReadTimeout    int    `yaml:"readTimeout" env-default:"20"`
		WriteTimeout   int    `yaml:"writeTimeout" env-default:"20"`
	} `yaml:"http" env-required:"true"`
	// DB represents configuration for database.
	DB struct {
		URL        string `env:"MONGO_URL" env-required:"true"`
		Database   string `yaml:"database" env-required:"true"`
		Collection string `yaml:"collection" env-required:"true"`
	} `yaml:"mongo" env-required:"true"`
}

var instance *Config
var once sync.Once

// Get loads .env file and config from given path.
// Returns config instance.
func Get(configPath string, dotenvPath string) *Config {
	logger := logger.GetLogger()

	logger.Info("loading .env file")
	if err := godotenv.Load(dotenvPath); err != nil {
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
