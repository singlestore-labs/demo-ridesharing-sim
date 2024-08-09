package exporter

import (
	"context"
	"fmt"
	"log"
	"simulator/config"
	"simulator/models"
	"time"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"github.com/twmb/franz-go/pkg/sr"
)

var KafkaClient *kgo.Client
var SchemaRegistryClient *sr.Client

var Serde sr.Serde

func InitializeKafkaClient() {
	if config.Kafka.Broker == "" {
		return
	}
	seeds := []string{config.Kafka.Broker}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("simulator"),
	)
	if err != nil {
		log.Fatalf("unable to create kafka client: %v", err)
	}
	KafkaClient = cl

	rcl, err := sr.NewClient(sr.URLs(config.Kafka.SchemaRegistry))
	if err != nil {
		log.Fatalf("unable to create schema registry client: %v", err)
	}
	SchemaRegistryClient = rcl

	for _, topic := range []string{"ridesharing-sim.trips"} {
		CreateTopic(topic)
	}

	CreateAvroSchemas()
}

func CreateTopic(topic string) {
	req := kmsg.NewPtrCreateTopicsRequest()
	t := kmsg.NewCreateTopicsRequestTopic()
	t.Topic = topic
	t.NumPartitions = 1
	t.ReplicationFactor = 1
	req.Topics = append(req.Topics, t)

	res, err := req.RequestWith(context.Background(), KafkaClient)
	if err != nil {
		log.Fatalf("unable to create kafka topic: %v", err)
	}

	if err := kerr.ErrorForCode(res.Topics[0].ErrorCode); err != nil && err != kerr.TopicAlreadyExists {
		log.Fatalf("kafka topic creation failure: %v", err)
		return
	}
	fmt.Printf("kafka topic %s created successfully!\n", t.Topic)
}

func KafkaProduceTrip(trip models.Trip) {
	KafkaClient.Produce(
		context.Background(),
		&kgo.Record{
			Key:   []byte(trip.ID),
			Topic: "ridesharing-sim.trips",
			Value: Serde.MustEncode(trip),
		},
		func(r *kgo.Record, err error) {
			if err != nil {
				fmt.Printf("unable to produce: %v", err)
			}
		},
	)
}

func KafkaDebugTripConsumer() {
	KafkaClient.AddConsumeTopics("ridesharing-sim.trips")
	for {
		fs := KafkaClient.PollFetches(context.Background())
		fs.EachRecord(func(r *kgo.Record) {
			var trip models.Trip
			err := Serde.Decode(r.Value, &trip)
			if err != nil {
				fmt.Printf("unable to decode: %v", err)
			}
			fmt.Printf("received: %v\n", trip)
		})
		time.Sleep(1 * time.Second)
	}
}
