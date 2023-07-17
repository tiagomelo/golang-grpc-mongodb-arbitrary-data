include .env
export

# ==============================================================================
# Help

.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]\n"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==============================================================================
# Proto

.PHONY: proto
## proto: compiles .proto files
proto:
	@ rm -rf api/proto/gen/productcatalog
	@ mkdir -p api/proto/gen/productcatalog
	@ cd api/proto ; \
	protoc --go_out=gen/productcatalog --go_opt=paths=source_relative --go-grpc_out=gen/productcatalog --go-grpc_opt=paths=source_relative productcatalog.proto

# ==============================================================================
# Docker-compose

.PHONY: start-mongodb
## start-mongodb: starts mongodb instance used for the app
start-mongodb:
	@ docker-compose up mongodb -d
	@ echo "Waiting for MongoDB to start..."
	@ until docker exec $(MONGODB_DATABASE_CONTAINER_NAME) mongosh --eval "db.adminCommand('ping')" >/dev/null 2>&1; do \
		echo "MongoDB not ready, sleeping for 5 seconds..."; \
		sleep 5; \
	done
	@ echo "MongoDB is up and running."

.PHONY: stop-mongodb
## stop-mongodb: stops mongodb instance used for the app
stop-mongodb:
	@ docker-compose stop mongodb

.PHONY: start-test-mongodb
## start-test-mongodb: starts mongodb instance used for integration tests
start-test-mongodb:
	@ docker-compose up mongodb_test -d
	@ echo "Waiting for Test MongoDB to start..."
	@ until docker exec $(MONGODB_TEST_DATABASE_CONTAINER_NAME) mongosh --eval "db.adminCommand('ping')" >/dev/null 2>&1; do \
		echo "Test MongoDB not ready, sleeping for 5 seconds..."; \
		sleep 5; \
	done
	@ echo "Test MongoDB is up and running."

.PHONY: stop-test-mongodb
## stop-test-mongodb: stops mongodb instance used for integration tests
stop-test-mongodb:
	@ docker-compose stop mongodb_test

.PHONY: stop-all-mongodb
## stop-all-mongodb: stops all mongodb instances
stop-all-mongodb:
	@ docker-compose down

# ==============================================================================
# Tests

.PHONY: test
## test: runs both unit and integration tests
test: start-test-mongodb
	@ go test -v ./...

# ==============================================================================
# Execution

.PHONY: run
## run: runs the gRPC server
run: start-mongodb
	@ go run cmd/main.go