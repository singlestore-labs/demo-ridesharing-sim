-- Create the database
DROP DATABASE IF EXISTS rideshare_demo;
CREATE DATABASE rideshare_demo;
USE rideshare_demo;

-- Create the trips table
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
    fare INT,
    distance DOUBLE,
    pickup_lat DOUBLE,
    pickup_long DOUBLE,
    dropoff_lat DOUBLE,
    dropoff_long DOUBLE,
    city VARCHAR(255),
    PRIMARY KEY (id)
);

-- Create the trips pipeline
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

-- Test the trips pipeline
TEST PIPELINE rideshare_kafka_trips LIMIT 1;

-- Start the trips pipeline
START PIPELINE rideshare_kafka_trips;

-- Create the riders table
DROP TABLE IF EXISTS riders;
CREATE TABLE riders (
    id VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone_number VARCHAR(255),
    date_of_birth DATETIME(6),
    created_at DATETIME(6),
    location_city VARCHAR(255),
    location_lat DOUBLE,
    location_long DOUBLE,
    status VARCHAR(20),
    PRIMARY KEY (id)
);

-- Create the riders pipeline
CREATE OR REPLACE PIPELINE rideshare_kafka_riders AS
    LOAD DATA KAFKA 'cqrik6h4mu94dmoo2370.any.us-east-1.mpx.prd.cloud.redpanda.com:9092/ridesharing-sim-riders'
    CONFIG '{"sasl.username": "<username>",
         "sasl.mechanism": "SCRAM-SHA-256",
         "security.protocol": "SASL_SSL",
         "ssl.ca.location": "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"}'
    CREDENTIALS '{"sasl.password": "<password>"}'
    DISABLE OUT_OF_ORDER OPTIMIZATION
    INTO TABLE riders
    FORMAT JSON
    (
        id <- id,
        first_name <- first_name,
        last_name <- last_name,
        email <- email,
        phone_number <- phone_number,
        @date_of_birth <- date_of_birth,
        @created_at <- created_at,
        location_city <- location_city,
        location_lat <- location_lat,
        location_long <- location_long,
        status <- status
    )
    SET date_of_birth = STR_TO_DATE(@date_of_birth, '%Y-%m-%dT%H:%i:%s.%f'),
        created_at = STR_TO_DATE(@created_at, '%Y-%m-%dT%H:%i:%s.%f')
    ON DUPLICATE KEY UPDATE
        first_name = VALUES(first_name),
        last_name = VALUES(last_name),
        email = VALUES(email),
        phone_number = VALUES(phone_number),
        date_of_birth = VALUES(date_of_birth),
        created_at = VALUES(created_at),
        location_city = VALUES(location_city),
        location_lat = VALUES(location_lat),
        location_long = VALUES(location_long),
        status = VALUES(status)

-- Test the riders pipeline
TEST PIPELINE rideshare_kafka_riders LIMIT 1;

-- Start the riders pipeline
START PIPELINE rideshare_kafka_riders;

-- Create the drivers table
DROP TABLE IF EXISTS drivers;
CREATE TABLE drivers (
    id VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone_number VARCHAR(255),
    date_of_birth DATETIME(6),
    created_at DATETIME(6),
    location_city VARCHAR(255),
    location_lat DOUBLE,
    location_long DOUBLE,
    status VARCHAR(20),
    PRIMARY KEY (id)
);

-- Create the drivers pipeline
CREATE OR REPLACE PIPELINE rideshare_kafka_drivers AS
    LOAD DATA KAFKA 'cqrik6h4mu94dmoo2370.any.us-east-1.mpx.prd.cloud.redpanda.com:9092/ridesharing-sim-drivers'
    CONFIG '{"sasl.username": "<username>",
         "sasl.mechanism": "SCRAM-SHA-256",
         "security.protocol": "SASL_SSL",
         "ssl.ca.location": "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"}'
    CREDENTIALS '{"sasl.password": "<password>"}'
    DISABLE OUT_OF_ORDER OPTIMIZATION
    INTO TABLE drivers
    FORMAT JSON
    (
        id <- id,
        first_name <- first_name,
        last_name <- last_name,
        email <- email,
        phone_number <- phone_number,
        @date_of_birth <- date_of_birth,
        @created_at <- created_at,
        location_city <- location_city,
        location_lat <- location_lat,
        location_long <- location_long,
        status <- status
    )
    SET date_of_birth = STR_TO_DATE(@date_of_birth, '%Y-%m-%dT%H:%i:%s.%f'),
        created_at = STR_TO_DATE(@created_at, '%Y-%m-%dT%H:%i:%s.%f')
    ON DUPLICATE KEY UPDATE
        first_name = VALUES(first_name),
        last_name = VALUES(last_name),
        email = VALUES(email),
        phone_number = VALUES(phone_number),
        date_of_birth = VALUES(date_of_birth),
        created_at = VALUES(created_at),
        location_city = VALUES(location_city),
        location_lat = VALUES(location_lat),
        location_long = VALUES(location_long),
        status = VALUES(status)

-- Test the drivers pipeline
TEST PIPELINE rideshare_kafka_drivers LIMIT 1;

-- Start the drivers pipeline
START PIPELINE rideshare_kafka_drivers;