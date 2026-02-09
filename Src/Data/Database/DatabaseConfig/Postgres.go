package DatabaseConfig

import (
	"RideMarket-CleanWebApi-GoLang/Config"
	"RideMarket-CleanWebApi-GoLang/pkg/Logging/Log"
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbClient *gorm.DB
var logger = Log.NewLogger(Config.GetConfig())

func InitDb(cfg *Config.Config) error {
	var err error
	postgresConfig := cfg.Postgres
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		postgresConfig.Host, strconv.Itoa(postgresConfig.Port), postgresConfig.User, postgresConfig.Password, postgresConfig.DbName,
		postgresConfig.SslMode)

	dbClient, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDb, _ := dbClient.DB()
	err = sqlDb.Ping()
	if err != nil {
		return err
	}

	sqlDb.SetConnMaxIdleTime(postgresConfig.ConnMaxIdleTime)
	sqlDb.SetMaxOpenConns(postgresConfig.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(postgresConfig.ConnMaxLifetime)

	logger.Info(Log.Postgres, Log.Startup, "Connection Is Successfully Established", nil)

	err = MigrateAndSeed()
	if err != nil {
		logger.Info(Log.Postgres, Log.Migration, err.Error(), nil)

	}

	return nil
}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	connection, err := dbClient.DB()
	if err != nil {
		log.Println(err)
	}
	err = connection.Close()
	if err != nil {
		return
	}

}
