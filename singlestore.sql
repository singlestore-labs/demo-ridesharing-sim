DROP DATABASE IF EXISTS rideshare_demo;
CREATE DATABASE rideshare_demo;

USE rideshare_demo;

DROP PIPELINE IF EXISTS rideshare_kafka_trips;
DROP TABLE IF EXISTS trips;
CREATE TABLE trips (
    id VARCHAR(255) NOT NULL,
    driver_id VARCHAR(255),
    rider_id VARCHAR(255),
    status VARCHAR(20),
    request_time DATETIME(6),
    accept_time DATETIME(6),
    pickup_time DATETIME(6),
    dropoff_time DATETIME(6),
    fare INT NOT NULL,
    distance DOUBLE NOT NULL,
    pickup_lat DOUBLE NOT NULL,
    pickup_long DOUBLE NOT NULL,
    dropoff_lat DOUBLE NOT NULL,
    dropoff_long DOUBLE NOT NULL,
    city VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

CREATE OR REPLACE PIPELINE rideshare_kafka_trips AS
    LOAD DATA KAFKA 'cqrik6h4mu94dmoo2370.any.us-east-1.mpx.prd.cloud.redpanda.com:9092/ridesharing-sim-trips'
    CONFIG '{"sasl.username": "<username>",
         "sasl.mechanism": "SCRAM-SHA-256",
         "security.protocol": "SASL_SSL",
         "ssl.ca.location": "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"}'
    CREDENTIALS '{"sasl.password": "<password>"}'
    DISABLE OUT_OF_ORDER OPTIMIZATION
    INTO TABLE trips
    FORMAT JSON
    (
        id <- id,
        rider_id <- rider_id,
        driver_id <- driver_id,
        status <- status,
        @request_time <- request_time,
        @accept_time <- accept_time,
        @pickup_time <- pickup_time,
        @dropoff_time <- dropoff_time,
        fare <- fare,
        distance <- distance,
        pickup_lat <- pickup_lat,
        pickup_long <- pickup_long,
        dropoff_lat <- dropoff_lat,
        dropoff_long <- dropoff_long,
        city <- city
    )
    SET request_time = STR_TO_DATE(@request_time, '%Y-%m-%dT%H:%i:%s.%f'),
    accept_time = STR_TO_DATE(@accept_time, '%Y-%m-%dT%H:%i:%s.%f'),
    pickup_time = STR_TO_DATE(@pickup_time, '%Y-%m-%dT%H:%i:%s.%f'),
    dropoff_time = STR_TO_DATE(@dropoff_time, '%Y-%m-%dT%H:%i:%s.%f')
    ON DUPLICATE KEY UPDATE
        rider_id = VALUES(rider_id),
        driver_id = VALUES(driver_id),
        status = VALUES(status),
        request_time = VALUES(request_time),
        accept_time = VALUES(accept_time),
        pickup_time = VALUES(pickup_time),
        dropoff_time = VALUES(dropoff_time),
        fare = VALUES(fare),
        distance = VALUES(distance),
        pickup_lat = VALUES(pickup_lat),
        pickup_long = VALUES(pickup_long),
        dropoff_lat = VALUES(dropoff_lat),
        dropoff_long = VALUES(dropoff_long),
        city = VALUES(city);

TEST PIPELINE rideshare_kafka_trips LIMIT 1;

START PIPELINE rideshare_kafka_trips;

SELECT status, COUNT(*) as trip_count
    FROM trips
    GROUP BY status
    ORDER BY status;