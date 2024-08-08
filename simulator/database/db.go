package database

import "simulator/config"

func Initialize() {
	InitializeLocal()
	if config.Store == "singlestore" {
		connectSingleStore()
	} else if config.Store == "snowflake" {
		connectSnowflake()
	}
}
