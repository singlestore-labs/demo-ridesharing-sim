KAFKA_TOPIC="ridesharing-sim.trips"
eval $(cat .env)

URL="https://$SNOWFLAKE_ACCOUNT.snowflakecomputing.com"
NAME="ridesharing-sim.trips"
DB_NAME="RIDESHARE_INGEST"
TABLE_NAME="TRIPS"

curl -i -X PUT -H "Content-Type:application/json" \
    "http://localhost:8083/connectors/$NAME/config" \
    -d '{
        "connector.class":"com.snowflake.kafka.connector.SnowflakeSinkConnector",
        "errors.log.enable":"true",
        "snowflake.database.name":"'$DB_NAME'",
        "snowflake.private.key":"'$PRIVATE_KEY'",
        "snowflake.schema.name":"PUBLIC",
        "snowflake.role.name":"'$DB_NAME'",
        "snowflake.url.name":"'$URL'",
        "snowflake.user.name":"'$SNOWFLAKE_USER'",
        "snowflake.enable.schematization": "TRUE",
        "snowflake.ingestion.method": "SNOWPIPE_STREAMING",
        "topics":"'$KAFKA_TOPIC'",
        "name":"'$NAME'",
        "key.converter":"org.apache.kafka.connect.storage.StringConverter",
        "value.converter":"io.confluent.connect.avro.AvroConverter",
        "value.converter.schema.registry.url": "http://host.docker.internal:18081",
        "buffer.count.records":"1000000",
        "buffer.flush.time":"10",
        "buffer.size.bytes":"250000000",
        "snowflake.topic2table.map":"'$KAFKA_TOPIC:$TABLE_NAME'"
    }'