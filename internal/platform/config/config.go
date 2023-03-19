package config

import (
	"github.com/richmondgoh8/go-casbin/pkg/client/postgres"
	"github.com/simukti/sqldb-logger"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB *postgres.DBConfig
}

func InitReader() {
	environment := ""
	if len(os.Args) < 2 {
		log.Fatalf("Env not supplied in argument")
	} else {
		environment = os.Args[1]
	}

	err := godotenv.Load(environment + ".env")
	if err != nil {
		log.Fatalf("Error loading %s.env file", environment)
	}
}

func Init() AppConfig {
	var minLevel uint64 = 1
	minLevelStr := os.Getenv("DB_LOGGER_LEVEL")
	if minLevelStr != "" {
		if newMinLevel, err := strconv.ParseUint(minLevelStr, 10, 64); err == nil {
			minLevel = newMinLevel
		}
	}

	appConfig := AppConfig{
		DB: &postgres.DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Schema:   os.Getenv("DB_SCHEMA"),
			MinLevel: sqldblogger.Level(minLevel),
		},
	}

	return appConfig
}
