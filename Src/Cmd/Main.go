package Cmd

import (
	"RideMarket-CleanWebApi-GoLang/Api"
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/Data/Cache"
)

func main() {
	cfg := Config.GetConfig()

	Cache.InitRedis(cfg)
	defer Cache.CloseRedis()

	Api.InitServer(cfg)
}
