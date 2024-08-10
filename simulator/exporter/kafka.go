package exporter

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"simulator/config"
	"simulator/models"
	"time"

	"github.com/goccy/go-json"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"github.com/twmb/franz-go/pkg/sasl/scram"
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

	if config.Kafka.SASLUsername != "" && config.Kafka.SASLPassword != "" {
		tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
		opts := []kgo.Opt{
			kgo.SeedBrokers(seeds...),
			kgo.ConsumerGroup("simulator"),
			// kgo.SASL(plain.Auth{
			// 	User: config.Kafka.SASLUsername,
			// 	Pass: config.Kafka.SASLPassword,
			// }.AsMechanism()),
			kgo.SASL(scram.Auth{
				User: config.Kafka.SASLUsername,
				Pass: config.Kafka.SASLPassword,
			}.AsSha256Mechanism()),
			kgo.Dialer(tlsDialer.DialContext),
		}
		cl, err := kgo.NewClient(opts...)
		if err != nil {
			log.Fatalf("unable to create kafka client: %v", err)
		}
		KafkaClient = cl
	} else {
		cl, err := kgo.NewClient(
			kgo.SeedBrokers(seeds...),
			kgo.ConsumerGroup("simulator"),
		)
		if err != nil {
			log.Fatalf("unable to create kafka client: %v", err)
		}
		KafkaClient = cl
	}
	log.Printf("kafka client connected successfully!\n")

	rcl, err := sr.NewClient(sr.URLs(config.Kafka.SchemaRegistry))
	if err != nil {
		log.Fatalf("unable to create schema registry client: %v", err)
	}
	SchemaRegistryClient = rcl

	for _, topic := range []string{"ridesharing-sim-trips"} {
		CreateTopic(topic)
	}

	// CreateAvroSchemas()
}

func CreateTopic(topic string) {
	req := kmsg.NewPtrCreateTopicsRequest()
	t := kmsg.NewCreateTopicsRequestTopic()
	t.Topic = topic
	t.NumPartitions = 1
	t.ReplicationFactor = 3
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
	jsonTrip, err := json.Marshal(trip)
	if err != nil {
		log.Fatalf("unable to marshal trip: %v", err)
	}
	KafkaClient.Produce(
		context.Background(),
		&kgo.Record{
			Key:   []byte(trip.ID),
			Topic: "ridesharing-sim-trips",
			// Value: Serde.MustEncode(trip),
			Value: jsonTrip,
		},
		func(r *kgo.Record, err error) {
			if err != nil {
				fmt.Printf("unable to produce: %v", err)
			}
		},
	)
}

func KafkaDebugTripConsumer() {
	KafkaClient.AddConsumeTopics("ridesharing-sim-trips")
	for {
		fs := KafkaClient.PollFetches(context.Background())
		fs.EachRecord(func(r *kgo.Record) {
			var trip models.Trip
			// err := Serde.Decode(r.Value, &trip)
			err := json.Unmarshal(r.Value, &trip)
			if err != nil {
				fmt.Printf("unable to decode: %v", err)
			}
			fmt.Printf("received: %v\n", trip)
		})
		time.Sleep(1 * time.Second)
	}
}
