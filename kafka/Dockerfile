FROM docker.redpanda.com/redpandadata/connectors:v1.0.27

USER root

RUN mkdir -p /opt/kafka/redpanda-plugins/snowflake
RUN curl -o /opt/kafka/redpanda-plugins/snowflake/snowflake-kafka-connector-2.1.2.jar https://repo1.maven.org/maven2/com/snowflake/snowflake-kafka-connector/2.2.2/snowflake-kafka-connector-2.2.2.jar
RUN curl -o /opt/kafka/redpanda-plugins/snowflake/bc-fips-1.0.1.jar https://repo1.maven.org/maven2/org/bouncycastle/bc-fips/1.0.1/bc-fips-1.0.1.jar
RUN curl -o /opt/kafka/redpanda-plugins/snowflake/bcpkix-fips-1.0.3.jar https://repo1.maven.org/maven2/org/bouncycastle/bcpkix-fips/1.0.3/bcpkix-fips-1.0.3.jar

USER redpanda