package Database

import (
	"RideMarket-CleanWebApi-GoLang/Config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbClient *gorm.DB

func InitDb(cfg *Config.Config) error {
	postgresConfig := cfg.Postgres
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		postgresConfig.Host, postgresConfig.Port, postgresConfig.User, postgresConfig.Password, postgresConfig.DbName,
		postgresConfig.SslMode)

	dbClient, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
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

	log.Println("Connection Is Successfully Established")

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
	connection.Close()

}
