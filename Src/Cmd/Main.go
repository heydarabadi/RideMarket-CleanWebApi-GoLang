package main

import (
	"RideMarket-CleanWebApi-GoLang/Api"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Data/Cache"
	"RideMarket-CleanWebApi-GoLang/Data/Database"
	"context"
	"log"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg := Config.GetConfig()

	if err := Cache.InitRedis(cfg, ctx); err != nil {
		return err
	}
	defer Cache.CloseRedis()

	if err := Database.InitDb(cfg); err != nil {
		return err
	}
	defer Database.CloseDb()

	Api.InitServer(cfg)
	return nil
}
