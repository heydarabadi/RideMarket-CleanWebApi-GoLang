package Config

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Cors     CorsConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port    string
	RunMode string
}

type LoggerConfig struct {
	FilePath string
	Encoding string
	Level    string
}

type CorsConfig struct {
	AllowedOrigins string
}

type PostgresConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DbName          string
	SslMode         string
	ConnMaxIdleTime time.Duration
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host               string
	Port               int
	Password           string
	Db                 int
	ReadTimeOut        time.Duration
	PoolSize           int
	PoolTimeOut        time.Duration
	WriteTimeOut       time.Duration
	DialTimeOut        time.Duration
	IdleCheckFrequency time.Duration
}

func GetConfig() *Config {
	cfgPath := getConfigPath(os.Getenv("APP_ENV"))
	v, err := loadConfig(cfgPath, "yml")
	if err != nil {
		log.Fatalf("Load config err: %v", err)
	}
	cfg, err := parseConfig(v)
	if err != nil {
		log.Fatalf("Parse config err: %v", err)
	}
	return cfg
}

func parseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Unable To Parse Config %v", err)
		return nil, err
	}
	return &cfg, err
}

func loadConfig(fileName string, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(fileName)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable To Read Config %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		} else {
			return nil, err
		}
	}
	return v, nil

}

func getConfigPath(env string) string {
	if strings.ToLower(env) == "docker" {
		return "config/config-docker"
	} else if strings.ToLower(env) == "production" {
		return "config/config-production"
	} else if strings.ToLower(env) == "development" {
		return "../config/config-development"
	} else {
		return "../config/config-development"
	}
}
