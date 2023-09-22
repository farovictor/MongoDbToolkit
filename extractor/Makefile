SHELL := /bin/bash

include .env

export MONGO_CONN_URI
export MONGO_DBNAME
export MONGO_COLLECTION
export APPNAME

build:
	@go build -o ./bin/extractor -v ./src
	@chmod a+x ./bin/extractor

run-ping: build
	@./bin/extractor ping --conn-uri "$(MONGO_CONN_URI)"

run-collection-exists: build
	@./bin/extractor collxst --conn-uri "$(MONGO_CONN_URI)" --db-name "$(MONGO_DBNAME)" --collection "$(MONGO_COLLECTION)" --app-name "$(APPNAME)"

run-extraction: build
	@./bin/extractor extract --conn-uri "$(MONGO_CONN_URI)" --db-name "$(MONGO_DBNAME)" --collection "$(MONGO_COLLECTION)" --app-name "$(APPNAME)" --mapping record --query '{"latitude":{"$$gte":30}}'

run-extract-batch: build
	@./bin/extractor extract-batch \
		--conn-uri "$(MONGO_CONN_URI)" \
		--db-name "$(MONGO_DBNAME)" \
		--collection "$(MONGO_COLLECTION)" \
		--app-name "$(APPNAME)" \
		--mapping record \
		--query '{"latitude":{"$$gte":30}}' \
		--output-path "./data" \
		--output-prefix "mbl" \
		--num-concurrent-files 10

run-test:
	@go test ./...

.PHONY: run-ping, build, echo-path, run-test, run-collection--exists, run-extraction