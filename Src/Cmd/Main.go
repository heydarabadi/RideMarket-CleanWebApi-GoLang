package main

import (
	"RideMarket-CleanWebApi-GoLang/Api"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Data/Cache"
	"RideMarket-CleanWebApi-GoLang/Data/Database/DatabaseConfig"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"context"
	"time"
)

func main() {
	cfg := Config.GetConfig()

	logger := Log.NewLogger(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := Cache.InitRedis(cfg, ctx); err != nil {
		logger.Fatal(Log.Redis, Log.Startup, err.Error(), nil)
	}
	defer Cache.CloseRedis()

	if err := DatabaseConfig.InitDb(cfg); err != nil {
		logger.Fatal(Log.Postgres, Log.Startup, err.Error(), nil)

	}
	defer DatabaseConfig.CloseDb()

	Api.InitServer(cfg)
}
