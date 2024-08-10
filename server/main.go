package main

import (
	"fmt"
	"server/config"
	"server/database"
)

func main() {
	config.Verify()
	database.Initialize()
	fmt.Println("Hello, World!")
}
