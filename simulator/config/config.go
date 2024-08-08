package config

import (
	"os"

	"github.com/jaswdr/faker/v2"
)

var Port = "8080"
var NumRiders = 100
var NumDrivers = 70

var Faker = faker.New()

var Store = os.Getenv("STORE")

var SingleStore = struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}{
	os.Getenv("SINGLESTORE_HOST"),
	os.Getenv("SINGLESTORE_PORT"),
	os.Getenv("SINGLESTORE_DATABASE"),
	os.Getenv("SINGLESTORE_USER"),
	os.Getenv("SINGLESTORE_PASSWORD"),
}

var Snowflake = struct {
	Account   string
	User      string
	Password  string
	Warehouse string
	Database  string
	Schema    string
}{
	os.Getenv("SNOWFLAKE_ACCOUNT"),
	os.Getenv("SNOWFLAKE_USER"),
	os.Getenv("SNOWFLAKE_PASSWORD"),
	os.Getenv("SNOWFLAKE_WAREHOUSE"),
	os.Getenv("SNOWFLAKE_DATABASE"),
	os.Getenv("SNOWFLAKE_SCHEMA"),
}
