package database

import (
	"fmt"
	"log"
	"server/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SingleStoreDB *gorm.DB

func connectSingleStore() {
	if config.SingleStore.Host == "" || config.SingleStore.Port == "" || config.SingleStore.Username == "" || config.SingleStore.Password == "" || config.SingleStore.Database == "" {
		log.Println("SingleStore configuration is not set")
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", config.SingleStore.Username, config.SingleStore.Password, config.SingleStore.Host, config.SingleStore.Port, config.SingleStore.Database)

	dialector := mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
		DontSupportForShareClause: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to singlestore database: %v", err)
	}
	SingleStoreDB = db
	log.Println("Successfully connected to SingleStore")
}
