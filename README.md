# Ridesharing Simulation

**Attention**: The code in this repository is intended for experimental use only and is not fully tested, documented, or supported by SingleStore. Visit the [SingleStore Forums](https://www.singlestore.com/forum/) to ask questions about this repository.

## Usage

1. [Sign up](https://www.singlestore.com/cloud-trial/) for the SingleStore Free Shared Tier.
2. Create a public/private key pair for the kafka connector to use.
```
openssl genrsa 4096 | openssl pkcs8 -topk8 -inform PEM -out kafka_key.p8 -nocrypt
openssl rsa -in kafka_key.p8 -pubout -out kafka_key.pub
```

## Resources

* [Documentation](https://docs.singlestore.com)
* [Twitter](https://twitter.com/SingleStoreDevs)
* [SingleStore Forums](https://www.singlestore.com/forum)

