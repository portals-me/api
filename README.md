# portals-me

## JWT

Server will generate the own jwt for authentication and authorization.

### Key Generation

```sh
$ cd token
$ ssh-keygen -t ecdsa -b 256 -f jwtES256.key
$ openssl ec -in jwtES256.key -pubout -outform PEM -out jwtES256.key.pub
```

## Tests

cf: [https://medium.com/@octoz/automate-your-serverless-api-integration-tests-locally-e2f41d3ec757](https://medium.com/@octoz/automate-your-serverless-api-integration-tests-locally-e2f41d3ec757)
