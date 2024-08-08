package database

import (
	"database/sql"
	"fmt"
	"simulator/config"

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
		fmt.Printf("Failed to create DSN: %v", err)
	}
	// Connect to Snowflake
	db, err = sql.Open("snowflake", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to Snowflake: %v", err)
	}
	// Test the connection
	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to ping Snowflake: %v", err)
	}
	fmt.Println("Successfully connected to Snowflake")
}
