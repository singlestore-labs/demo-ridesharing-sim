package config

import (
	"os"
	"slices"
	"strconv"

	"github.com/jaswdr/faker/v2"
)

var Port = "8080"

var numRiders = os.Getenv("NUM_RIDERS")
var numDrivers = os.Getenv("NUM_DRIVERS")

var NumRiders = 100
var NumDrivers = 70

var City = os.Getenv("CITY")
var validCities = []string{"San Francisco"}

var Faker = faker.New()

var Kafka = struct {
	Broker         string
	SchemaRegistry string
}{
	Broker:         os.Getenv("KAFKA_BROKER"),
	SchemaRegistry: os.Getenv("KAFKA_SCHEMA_REGISTRY"),
}

func Verify() {
	if Port == "" {
		Port = "8001"
	}
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
