package exporter

import (
	"context"
	"fmt"
	"log"
	"simulator/models"

	"github.com/hamba/avro"
	"github.com/twmb/franz-go/pkg/sr"
)

var tripSchema = avro.MustParse(`
{
  "type": "record",
  "name": "trip",
  "fields": [
    {"name": "id", "type": "string"},
    {"name": "driver_id", "type": "string"},
    {"name": "rider_id", "type": "string"},
    {"name": "status", "type": "string"},
    {"name": "request_time", "type": {"type": "long", "logicalType": "timestamp-millis"}},
    {"name": "accept_time", "type": {"type": "long", "logicalType": "timestamp-millis"}},
    {"name": "pickup_time", "type": {"type": "long", "logicalType": "timestamp-millis"}},
    {"name": "dropoff_time", "type": {"type": "long", "logicalType": "timestamp-millis"}},
    {"name": "fare", "type": "int"},
    {"name": "distance", "type": "double"},
    {"name": "pickup_lat", "type": "double"},
    {"name": "pickup_long", "type": "double"},
    {"name": "dropoff_lat", "type": "double"},
    {"name": "dropoff_long", "type": "double"},
    {"name": "city", "type": "string"}
  ]
}
`)

func CreateAvroSchemas() {
	// Create the trip schema
	topic := "ridesharing-sim-trips"
	ss, err := SchemaRegistryClient.CreateSchema(context.Background(), fmt.Sprintf("%s-value", topic), sr.Schema{
		Schema: tripSchema.String(),
		Type:   sr.TypeAvro,
	})
	if err != nil {
		log.Fatalf("unable to create schema: %v", err)
	}
	fmt.Printf("created or reusing schema subject %q version %d id %d\n", ss.Subject, ss.Version, ss.ID)

	Serde.Register(
		ss.ID,
		models.Trip{},
		sr.EncodeFn(func(v any) ([]byte, error) {
			return avro.Marshal(tripSchema, v)
		}),
		sr.DecodeFn(func(b []byte, v any) error {
			return avro.Unmarshal(tripSchema, b, v)
		}),
	)
}
