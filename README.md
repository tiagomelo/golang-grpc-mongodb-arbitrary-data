# golang-grpc-mongodb-arbitrary-data

This is an example of a CRUD GRPC API using [MongoDB](https://mongodb.com/). It shows how we can handle arbitrary data types.

## running it

```
$ make run
```

## testing it

Both unit and integration tests are provided.

```
$ make test
```

## Available Makefile targets

```
$ make help

Usage: make [target]

  help                 shows this help message
  proto                compiles .proto files
  start-mongodb        starts mongodb instance used for the app
  stop-mongodb         stops mongodb instance used for the app
  start-test-mongodb   starts mongodb instance used for integration tests
  stop-test-mongodb    stops mongodb instance used for integration tests
  stop-all-mongodb     stops all mongodb instances
  test                 runs both unit and integration tests
  run                  runs the gRPC server
```