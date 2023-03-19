package postgres

import (
	"fmt"
	"github.com/richmondgoh8/go-casbin/pkg/middleware/logger"
	"log"
	"os"
	"strconv"
	"time"

	// postgres db driver
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"

	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zapadapter"
)

const DriverName = "postgres"

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Schema   string
	MinLevel sqldblogger.Level
}

func Init(dbConfig *DBConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Name)

	client, err := sqlx.Open(DriverName, dataSource)
	if err != nil {
		return nil, err
	}
	logger.InitBasic()
	loggerAdapter := zapadapter.New(logger.ZapBasicLogger)
	client = sqlx.NewDb(sqldblogger.OpenDriver(dataSource, client.Driver(), loggerAdapter, sqldblogger.WithMinimumLevel(dbConfig.MinLevel)), DriverName)

	connectionMaxLifeTime, err := strconv.Atoi(os.Getenv("POSTGRES_CONNECTION_MAX_LIFETIME"))
	if err != nil {
		log.Println("Unable to get POSTGRES_CONNECTION_MAX_LIFETIME from config. Setting default value of 0")
		connectionMaxLifeTime = 0
	}
	maxOpenConnections, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNECTIONS"))
	if err != nil {
		log.Println("Unable to get POSTGRES_MAX_OPEN_CONNECTIONS from config. Setting default value of 10")
		maxOpenConnections = 10
	}
	maxIdleConnections, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNECTIONS"))
	if err != nil {
		log.Println("Unable to get POSTGRES_MAX_IDLE_CONNECTIONS from config. Setting default value of 10")
		maxIdleConnections = 10
	}

	client.SetConnMaxLifetime(time.Duration(connectionMaxLifeTime))
	client.SetMaxOpenConns(maxOpenConnections)
	client.SetMaxIdleConns(maxIdleConnections)

	// verifies connection is db is working
	if err := client.Ping(); err != nil {
		return nil, err
	}

	return client, nil
}
