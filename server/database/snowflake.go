package database

import (
	"database/sql"
	"log"
	"server/config"

	sf "github.com/snowflakedb/gosnowflake"
)

var db *sql.DB

func connectSnowflake() {
	cfg := &sf.Config{
		Account:   config.Snowflake.Account,
		User:      config.Snowflake.User,
		Password:  config.Snowflake.Password,
		Database:  config.Snowflake.Database,
		Schema:    config.Snowflake.Schema,
		Warehouse: config.Snowflake.Warehouse,
	}
	dsn, err := sf.DSN(cfg)
	if err != nil {
		log.Fatalf("Failed to create DSN: %v", err)
	}
	// Connect to Snowflake
	db, err = sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to Snowflake: %v", err)
	}
	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping Snowflake: %v", err)
	}
	log.Println("Successfully connected to Snowflake")
}
