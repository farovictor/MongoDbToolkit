SHELL := /bin/bash

include .env

export MONGO_CONN_URI
export MONGO_DBNAME
export MONGO_COLLECTION
export APPNAME

build:
	@go build -o ./bin/loader -v ./src
	@chmod a+x ./bin/loader

build-container:
	@export GIT_COMMIT=$$(git rev-list -1 HEAD); export BUILD_TIME=$$(date -u +'%Y-%m-%dT%H:%M:%SZ'); export VERSION=1.0.0; \
	docker build \
				--build-arg GIT_COMMIT=$$GIT_COMMIT \
				--build-arg BUILD_TIME=$$BUILD_TIME \
				--build-arg VERSION=$$VERSION \
				-t "mongodb-loader:test" .

run-ping: build
	@./bin/loader --log-level=debug ping --conn-uri "$(MONGO_CONN_URI)" --db-name "$(MONGO_DBNAME)" --app-name "$(APPNAME)"

run-collection-exists: build
	@./bin/loader collxst --conn-uri "$(MONGO_CONN_URI)" --db-name "$(MONGO_DBNAME)" \
		--collection "$(MONGO_COLLECTION)" --app-name "$(APPNAME)"

run-load: build
	./bin/loader load --conn-uri "$(MONGO_CONN_URI)" --db-name "$(MONGO_DBNAME)" \
		--collection "$(MONGO_COLLECTION)" --app-name "$(APPNAME)" \
		--file-prefix "mbl_74c0476e-5be2-11ee-a791-96e281524756" --search-path "./data"

run-load-batch: build
	@./bin/loader load-batch \
		--log-level "info" \
		--conn-uri "$(MONGO_CONN_URI)" \
		--db-name "$(MONGO_DBNAME)" \
		--collection "$(MONGO_COLLECTION)" \
		--app-name "$(APPNAME)" \
		--search-path "./data" \
		--file-prefix "mbl_" \
		--num-concurrent-files 10

run-test:
	@go test ./...

.PHONY: run-ping, build, run-test, run-collection-exists, run-load, run-load-batch