package Log

import (
	"RideMarket-CleanWebApi-GoLang/Config"
)

type ILogger interface {
	Init()

	Info(cat Category, subCat SubCategory, message string,
		extra map[Extrakey]interface{})

	Infof(template string, args ...interface{})

	Debug(cat Category, subCat SubCategory, message string,
		extra map[Extrakey]interface{})

	Debugf(template string, args ...interface{})

	Warning(cat Category, subCat SubCategory, message string,
		extra map[Extrakey]interface{})

	Warningf(template string, args ...interface{})

	Error(cat Category, subCat SubCategory, message string,
		extra map[Extrakey]interface{})

	Errorf(template string, args ...interface{})

	Fatal(cat Category, subCat SubCategory, message string,
		extra map[Extrakey]interface{})

	Fatalf(template string, args ...interface{})
}

func NewLogger(config *Config.Config) ILogger {
	return newZapLogger(config)
}
