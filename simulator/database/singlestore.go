package database

import (
	"fmt"
	"simulator/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var singlestoreDB *gorm.DB

func connectSingleStore() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC", config.SingleStore.User, config.SingleStore.Password, config.SingleStore.Host, config.SingleStore.Port, config.SingleStore.Database)

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
		fmt.Printf("failed to connect to database: %v", err)
	}
	singlestoreDB = db
}
