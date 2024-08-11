package config

import (
	"os"
	"slices"
	"strconv"

	"github.com/jaswdr/faker/v2"
)

var numRiders = os.Getenv("NUM_RIDERS")
var numDrivers = os.Getenv("NUM_DRIVERS")

var NumRiders = 100
var NumDrivers = 70

var City = os.Getenv("CITY")
var validCities = []string{"San Francisco", "San Jose"}

var Faker = faker.New()

var Kafka = struct {
	Broker         string
	SchemaRegistry string
	SASLUsername   string
	SASLPassword   string
}{
	Broker:         os.Getenv("KAFKA_BROKER"),
	SchemaRegistry: os.Getenv("KAFKA_SCHEMA_REGISTRY"),
	SASLUsername:   os.Getenv("KAFKA_SASL_USERNAME"),
	SASLPassword:   os.Getenv("KAFKA_SASL_PASSWORD"),
}

func Verify() {
	if num, err := strconv.ParseInt(numRiders, 10, 64); err == nil {
		NumRiders = int(num)
	}
	if num, err := strconv.ParseInt(numDrivers, 10, 64); err == nil {
		NumDrivers = int(num)
	}
	if slices.Contains(validCities, City) {
		City = "San Francisco"
	}
}
