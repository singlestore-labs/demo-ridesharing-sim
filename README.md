# Ridesharing Simulation

**Attention**: The code in this repository is intended for experimental use only and is not fully tested, documented, or supported by SingleStore. Visit the [SingleStore Forums](https://www.singlestore.com/forum/) to ask questions about this repository.

## Overview

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/assets/overview_dark.png">
  <source media="(prefers-color-scheme: light)" srcset="/assets/overview_light.png">
  <img alt="Ridesharing demo architecture" src="/assets/overview_light.png">
</picture>



## Getting Started

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
5. Set the flush time to 1 second.
6. Create a new API Key to connect to the kafka broker.

### Create Snowflake tables

1. Run lines 32-231 in the Snowflake worksheet. This will create a `RIDERS`, `DRIVERS`, and `TRIPS` table, and sets up tasks to merge in data from the stage tables every minute. 
   - This is required because the default Snowflake Kafka connector only supports inserting data into Snowflake tables, not updating them.
   - Our simulator relies on upserting data into the `TRIPS` table to update the status of a trip, as well as updating the `RIDERS` and `DRIVERS` tables to send location and status updates.

### Populate sample trips

1. If you want to populate the `TRIPS` table with sample data, you can use this [trips.csv](https://bk1031.s3.us-west-2.amazonaws.com/rideshare/trips.csv) file.
2. There's a placeholder in the Snowflake worksheet (lines 238-259) to load a csv file from an S3 bucket into the `TRIPS` table.

### Run the simulator and server

1. Run `make build` to build the docker images.
2. Run `docker compose up` to start everything.
3. You should see kafka topics being created and trips being generated.

### Run the frontend

1. In the `web/` directory, run `npm install` to install the dependencies.
2. Run `npm run dev` to start the frontend.
3. Open your browser and navigate to `http://localhost:5173` to view the frontend.
4. You should now be able to see the trips being generated in the frontend.

### Create iceberg table

1. Run lines 261-274 to create an external volume on S3 for the iceberg table. Make sure to follow the AWS instructions to create the required roles and policies.
2. Then run lines 276-286 to create the iceberg table and copy over the data from the `TRIPS` table.

### SingleStore Setup

1. [Sign up](https://www.singlestore.com/cloud-trial/) for the SingleStore Free Shared Tier.
2. Create a public/private key pair for the kafka connector to use.

## Resources

* [Documentation](https://docs.singlestore.com)
* [Twitter](https://twitter.com/SingleStoreDevs)
* [SingleStore Forums](https://www.singlestore.com/forum)