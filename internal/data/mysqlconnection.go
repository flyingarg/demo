package data

import (
	"database/sql"
	"demo/internal/env"
	"demo/internal/logger"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func GetConnection() *sql.DB {
	if DB != nil {
		return DB
	}

	userName := env.GetEnv("DB_USERNAME").(string)
	password := env.GetEnv("DB_PASSWORD").(string)
	database := env.GetEnv("DB_DATABASE").(string)
	hostname := env.GetEnv("DB_HOSTNAME").(string)
	connectionURL := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", userName, password, hostname, database)
	db, err := sql.Open("mysql", connectionURL)
	if err != nil {
		logger.Log.Sugar().Errorf("failed to create connection to database %s, using credentials %s:%s, with error %e", database, userName, password, err)
	}
	DB = db
	return DB
}

func CloseConnection() {
	if DB != nil {
		logger.Log.Sugar().Info("closing database connection ...")
		defer DB.Close()
	}
}
