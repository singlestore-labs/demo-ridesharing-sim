# Ridesharing Simulation

**Attention**: The code in this repository is intended for experimental use only and is not fully tested, documented, or supported by SingleStore. Visit the [SingleStore Forums](https://www.singlestore.com/forum/) to ask questions about this repository.

<img src="/assets/demo.gif" alt="Ridesharing demo" width="50%">

## Overview

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/assets/overview_dark.png">
  <source media="(prefers-color-scheme: light)" srcset="/assets/overview_light.png">
  <img alt="Ridesharing demo architecture" src="/assets/overview_light.png">
</picture>

Ride-sharing apps such as Uber and Lyft generate massive amounts of data every day. Being able to efficiently ingest and analyze this data is key to making crucial data-driven decisions. This demo showcases how SingleStore can be used to accelerate an existing analytics dashboard, enabling low-latency analytics on real-time data.

This demo consists of three main components:
- [Simulator](#simulator)
- [API Server](#api-server)
- [React Dashboard](#react-dashboard)

Our simulator generates realistic ride-sharing trip data and streams it to a Kafka topic. Using the Snowflake Kafka Connector, this data is then ingested into Snowflake tables. An API Server queries this data and exposes it through a RESTful interface. Finally, a React Dashboard consumes this API to provide dynamic visualizations of rider, driver, and trip information.

We'll then demonstrate how to seamlessly ingest trip data from Snowflake into SingleStore using Iceberg tables, achieving zero-ETL ingestion with minimal overhead. By leveraging SingleStore Pipelines to consume updates from our existing Kafka topics, we enable real-time analytics on our dashboard, showcasing SingleStore's ability to handle high-throughput, low-latency data processing and querying.

## Getting Started

You can check out a video walkthrough of this demo [here](https://drive.google.com/file/d/1bOQM5RJ3KkL1Pz6--AWWOIofk4CI-SHd/view?usp=sharing).

### Prerequisites

- Docker
- Golang (optional)
- Node.js (optional)

### Snowflake setup

1. Create a new SQL worksheet on Snowflake and copy in the contents of `snowflake.sql`.
2. Run the first 20 lines in the worksheet to setup the required workspace, database, user, and roles.
3. Run `make keygen` to generate the public key for the kafka connector. You should see output like this:
    ```
    ALTER USER RIDESHARE_INGEST SET RSA_PUBLIC_KEY='';
    PRIVATE_KEY=
    ```
4. Copy the output from the `ALTER USER` line and replace line 22 in `snowflake.sql`.
5. Copy the output from the `PRIVATE_KEY` line and hold onto it for the next step.

### Kafka setup

1. For this demo, we're using [Confluent Cloud](https://www.confluent.io/) as our kafka broker. You can sign up for a free trial [here](https://www.confluent.io/get-started/).
2. Create a Snowflake connector using the private key you generated in the previous step.
3. Set the topic to table mapping to `ridesharing-sim-trips:TRIPS_STAGE,ridesharing-sim-riders:RIDERS_STAGE,ridesharing-sim-drivers:DRIVERS_STAGE`.
4. Set the ingestion method to `SNOWPIPE_STREAMING` and the input format to `JSON`.
5. Make sure that "enable schematization" and "include createtime in metadata" are set to true.
6. Set the flush time to 1 second.
7. Create a new API Key to connect to the kafka broker.

### Verify the kafka connector

1. The kafka connector should have created the tables `RIDERS_STAGE`, `DRIVERS_STAGE`, and `TRIPS_STAGE`.
2. Run lines 27-30 in the Snowflake worksheet to verify that the kafka connector created the tables.

### Create Snowflake tables

1. Now that our staging tables are created, we are going to create the tables where our final data will land in.
2. Run lines 32-231 in the Snowflake worksheet. This will create a `RIDERS`, `DRIVERS`, and `TRIPS` table, and sets up Snowpipe tasks to merge in data from the stage tables every minute.
   - This is required because the default Snowflake Kafka connector only supports inserting data into Snowflake tables, not updating them.
   - Our simulator relies on upserting data into the `TRIPS` table to update the status of a trip, as well as updating the `RIDERS` and `DRIVERS` tables to send location and status updates.

### Populate sample trips

1. If you want to populate the `TRIPS` table with sample data, you can use this [trips.csv](https://bk1031.s3.us-west-2.amazonaws.com/rideshare/trips.csv) file.
2. There's a placeholder in the Snowflake worksheet (lines 238-259) to load a csv file from an S3 bucket into the `TRIPS` table.

### Run the simulator, server, and dashboard

1. Copy `example.env` to `.env` and modify the values to match your setup. You can ignore the SingleStore variables for now.
2. Run `make build` to build the docker images.
3. Run `docker compose up` to start everything.
4. You should see kafka topics being created and trips being generated.
5. The API server should be running on port 8000.
6. The React dashboard should be running on port 8080.
7. Open your browser and navigate to `http://localhost:8080` to view the frontend. You should be able to see the trips being generated.

> [!NOTE]
> The simulator currently uses `SASL_PLAIN` authentication to connect to the kafka broker. If you want to use `SCRAM`, you'll need to modify the code in `simulator/exporter/kafka.go`.

### Create iceberg table

1. Now that we have some data in our trips table, let's look at how to perform a zero-ETL ingest of this data into SingleStore without burning any Snowflake credits.
2. Run lines 261-274 to create an external volume on S3 for the iceberg table. Make sure to follow the AWS instructions to create the required roles and policies.
3. Then run lines 276-286 to create the iceberg table and copy over the data from the `TRIPS` table.
4. We should now have a `TRIPS_ICE` table with its catalog on S3.

### SingleStore Setup

1. [Sign up](https://www.singlestore.com/cloud-trial/) for the SingleStore Cloud Trial. You will need at least a S-2 workspace for this demo.
2. Create a new workspace and import the `singlestore.ipynb` notebook.
3. Edit the `CONFIG` and `CREDENTIALS` JSON in the 5th code cell so SingleStore can connect to your iceberg catalog.
4. Edit the `CONFIG` and `CREDENTIALS` JSON in cells 7, 10, and 13 so SingleStore can connect to your kafka broker.
5. Run the whole notebook.
6. Your tables should now be pulling in trip information from our kafka topics.
7. Update the SingleStore variables in your `.env` file, and restart the docker compose stack.
8. Open the React frontend, and select the SingleStore logo in the header. You should be able to see analytics from the trip data updating in real-time, as well as see the new pricing reccomendations in the analytics page.

## Simulator

The simulator has three main constructs:
- `Rider`: A person who requests a ride.
  - `id`: A unique identifier for the rider.
  - `first_name`: The first name of the rider.
  - `last_name`: The last name of the rider.
  - `email`: The email of the rider.
  - `phone_number`: The phone number of the rider.
  - `date_of_birth`: The date of birth of the rider.
  - `created_at`: The time the rider was created.
  - `location_city`: The city the rider is currently in.
  - `location_lat`: The latitude of the rider's current location.
  - `location_long`: The longitude of the rider's current location.
  - `status`: The status of the rider.
- `Driver`: A person who provides a ride.
  - `id`: A unique identifier for the driver.
  - `first_name`: The first name of the driver.
  - `last_name`: The last name of the driver.
  - `email`: The email of the driver.
  - `phone_number`: The phone number of the driver.
  - `date_of_birth`: The date of birth of the driver.
  - `created_at`: The time the driver was created.
  - `location_city`: The city the driver is currently in.
  - `location_lat`: The latitude of the driver's current location.
  - `location_long`: The longitude of the driver's current location.
  - `status`: The status of the driver.
- `Trip`: A ride given by a specific driver to a specific rider.
  - `id`: A unique identifier for the trip.
  - `driver_id`: The ID of the driver who is providing the ride.
  - `rider_id`: The ID of the rider who is requesting the ride.
  - `status`: The status of the trip.
  - `request_time`: The time the trip was requested.
  - `accept_time`: The time the trip was accepted.
  - `pickup_time`: The time the rider was picked up.
  - `dropoff_time`: The time the rider was dropped off.
  - `fare`: The fare of the trip.
  - `distance`: The distance of the trip.
  - `pickup_lat`: The latitude of the pickup location.
  - `pickup_long`: The longitude of the pickup location.
  - `dropoff_lat`: The latitude of the dropoff location.
  - `dropoff_long`: The longitude of the dropoff location.
  - `city`: The city the trip is in.

Using the provided configuration variables, the simulator will spawn in a number of riders and drivers at a random location in the specified city. Since there are a lot of cities in the bay area, we combined some of them together to create the following 16 cities.

<img alt="Bay area geojson map" src="/assets/map.png">

These cities are defined by GeoJSON polygons in the `simulator/data/` directory. Adding a new city is as simple as adding a new GeoJSON file, and updating the `ValidCities` array in `simulator/config/config.go`. Check out the `LoadGeoData()` function in `simulator/service/geo.go` to see how this array is used to load the polygons.

Each rider and driver run a loop in their own goroutines. The riders will request a ride to random location in the city, and wait for the ride to be accepted. This request creates a `Trip` object with the status set to `requested` and the `pickup_lat`, `pickup_long`, `dropoff_lat`, `dropoff_long` fields set to the rides's initial and requested coordinates. Drivers will accept the nearest ride request based on the pickup coordinates, setting the status to `accepted`. Then the driver proceeds to the rider's location by creating a straight line path between their current location and the rider's location, moving approximately 10 meters every 100 milliseconds. Once they have picked up the rider, the trips status is updated to `en_route`, and the driver continues to the dropoff location in the same manner. Once the rider is dropped off, the trip is marked as `completed`. The timestamps for each of these status updates is recorded in the appropriate fields.

The loop continues, and the rider and driver will continue to request and accept rides. Random delays are adding between each step in this loop to add some variance to the simulation.

At each step, the simulator will push the updated rider, driver, and trip objects to the following kafka topics:
- `ridesharing-sim-riders`: The rider objects.
- `ridesharing-sim-drivers`: The driver objects.
- `ridesharing-sim-trips`: The trip objects.

One thing to note is that riders and drivers are designed to be ephemeral, as they are generated on every run of the simulator. Their respective tables are only used to provide the real-time location and status of the rider and driver. They are not used to store historical data. The trips table is the only table that stores historical data. You may need to cleanup orphaned riders and drivers from time to time. You can do so with the following queries (works on both Snowflake and SingleStore):

```sql
DELETE FROM riders;
DELETE FROM drivers;
DELETE FROM trips WHERE status != 'completed';
```

### Configuration

Environment variables can be specified by creating a `.env` file in root directory. The following variables are supported:

- `NUM_RIDERS`: The number of riders to generate. (default: `100`)
- `NUM_DRIVERS`: The number of drivers to generate. (default: `70`)
- `CITY`: The city to generate trips in. (default: `San Francisco`)
- `KAFKA_BROKER`: The kafka broker to connect to.
- `KAFKA_SASL_USERNAME`: The username of the kafka broker.
- `KAFKA_SASL_PASSWORD`: The password of the kafka broker.

Note that `NUM_RIDERS`, `NUM_DRIVERS`, and `CITY` are set for each simulator instance in the docker compose file, and will not be pulled from the `.env` file.

### Running locally

If you want to run the simulator locally without docker, you can run `make run` in the `simulator/` directory. This will start a single instance of the simulator with the specified environment variables. Kafka topics will be created if they don't exist.

## API Server

The API server exposes a simple RESTful interface to make queries against the SingleStore and Snowflake databases. The database to query can be specified by the `db` query parameter (can be either `singlestore` or `snowflake`). Each endpoint will also return the query latency in microseconds through the `X-Query-Latency` header.

The following endpoints are available:
- `/trips/current/status`: Get the status breakdown of current active trips.
- `/trips/statistics`: Get total trip statistics.
- `/trips/statistics/daily`: Get daily trip statistics.
- `/trips/last/hour`: Get minute-by-minute trip counts for the last hour.
- `/trips/last/day`: Get hourly trip counts for the last day.
- `/trips/last/week`: Get daily trip counts for the last week.
- `/wait-time/last/hour`: Get minute-by-minute average wait times for the last hour.
- `/wait-time/last/day`: Get hourly average wait times for the last day.
- `/wait-time/last/week`: Get daily average wait times for the last week.
- `/riders`: Get current rider information.
- `/drivers`: Get current driver information.
- `/cities`: Get list of cities.

Not that the `city` query parameter can be used to filter results for a specific city in each of these requests.

Under the hood, the API server makes connections to both SingleStore and Snowflake, and runs the appropriate queries against the selected database. You can see the various queries defined for each database in the `server/service/` directory.

### Configuration

Environment variables can be specified by creating a `.env` file in root directory. The following variables are supported:

- `PORT`: The port the API server will run on. (default: `8000`)
- `SINGLESTORE_HOST`: The host of the SingleStore instance.
- `SINGLESTORE_PORT`: The port of the SingleStore instance.
- `SINGLESTORE_USERNAME`: The username of the SingleStore instance.
- `SINGLESTORE_PASSWORD`: The password of the SingleStore instance.
- `SINGLESTORE_DATABASE`: The database of the SingleStore instance.
- `SNOWFLAKE_ACCOUNT`: The account of the Snowflake instance.
- `SNOWFLAKE_USER`: The user of the Snowflake instance.
- `SNOWFLAKE_PASSWORD`: The password of the Snowflake instance.
- `SNOWFLAKE_WAREHOUSE`: The warehouse of the Snowflake instance.
- `SNOWFLAKE_DATABASE`: The database of the Snowflake instance.
- `SNOWFLAKE_SCHEMA`: The schema of the Snowflake instance.

If any of the SingleStore variables are blank, the API server will skip connecting to SingleStore. Similarly, if any of the Snowflake variables are blank, the API server will skip connecting to Snowflake.

### Running locally

If you want to run the API server locally without docker, you can run `make run` in the `server/` directory.

## React Dashboard

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/assets/dashboard_dark.png">
  <source media="(prefers-color-scheme: light)" srcset="/assets/dashboard_light.png">
  <img alt="React dashboard maps page" src="/assets/dashboard_light.png" width="49%">
</picture>
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/assets/analytics_dark.png">
  <source media="(prefers-color-scheme: light)" srcset="/assets/analytics_light.png">
  <img alt="React dashboard analytics page" src="/assets/analytics_light.png" width="49%">
</picture>

Our React dashboard consists of two main pages:
- A maps page that shows the current location of all rides and drivers on the map, as well as a summary of currently active trips.
- An analytics page that shows a variety of analytics on the ride data, such as the number of rides per hour and the average rider wait time.

The database being queried can be changed by selecting the appropriate logo in the header of the dashboard, and the refresh interval can be changed by selecting the appropriate option in the toolbar on the bottom right.

### Configuration

Environment variables can be specified by creating a `.env` file in the `web/` directory. The following variables are supported:

- `VITE_BACKEND_URL`: The URL of the API server. (default: `http://localhost:8000`)

### Running locally

If you want to run the dashboard locally, you can execute the following in the `web/` directory.

```
$ npm install
$ npm run dev
```

This will start the web server at `http://localhost:5173`. Make sure to set `VITE_BACKEND_URL` to point to the API server.

## Resources

* [Documentation](https://docs.singlestore.com)
* [Twitter](https://twitter.com/SingleStoreDevs)
* [SingleStore Forums](https://www.singlestore.com/forum)
