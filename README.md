# Ridesharing Simulation

**Attention**: The code in this repository is intended for experimental use only and is not fully tested, documented, or supported by SingleStore. Visit the [SingleStore Forums](https://www.singlestore.com/forum/) to ask questions about this repository.

## Overview

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/assets/overview_dark.png">
  <source media="(prefers-color-scheme: light)" srcset="/assets/overview_light.png">
  <img alt="Ridesharing demo architecture" src="/assets/overview_light.png">
</picture>



## Getting Started

There's a couple of different setup steps required to run this demo.

### Snowflake Setup

### Kafka Setup

### SingleStore Setup

1. [Sign up](https://www.singlestore.com/cloud-trial/) for the SingleStore Free Shared Tier.
2. Create a public/private key pair for the kafka connector to use.
```
openssl genrsa 4096 | openssl pkcs8 -topk8 -inform PEM -out kafka_key.p8 -nocrypt
openssl rsa -in kafka_key.p8 -pubout -out kafka_key.pub
cat kafka_key.pub | grep -v KEY- | tr -d '\012'
3. Copy the output and replace the placeholder public key in `snowflake.sql`
```

## Resources

* [Documentation](https://docs.singlestore.com)
* [Twitter](https://twitter.com/SingleStoreDevs)
* [SingleStore Forums](https://www.singlestore.com/forum)