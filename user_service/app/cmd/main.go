package main

import (
	"context"
	"flag"

	"github.com/juicyluv/sueta/user_service/app/config"
	"github.com/juicyluv/sueta/user_service/app/internal/user"
	"github.com/juicyluv/sueta/user_service/app/internal/user/db"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"github.com/juicyluv/sueta/user_service/app/pkg/mongo"
	"github.com/julienschmidt/httprouter"
)

var (
	configPath = flag.String("config-path", "app/config/config.yml", "path for application configuration file")
)

func main() {
	flag.Parse()
	logger.Init()

	logger := logger.GetLogger()
	logger.Info("logger initialized")

	cfg := config.Get(*configPath)
	logger.Info("loaded config file")

	router := httprouter.New()
	_ = router
	logger.Info("initialized httprouter")

	logger.Info("connecting to database")
	mongoClient, err := mongo.NewMongoClient(context.Background(),
		cfg.DB.Database, cfg.DB.URL)
	if err != nil {
		logger.Fatalf("cannot connect to mongodb: %v", err)
	}
	logger.Info("connected to database")

	userStorage := db.NewStorage(mongoClient, cfg.DB.Database)
	userService := user.NewService(userStorage, logger)
	_ = userService

}
