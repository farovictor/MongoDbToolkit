# Usage
This tool can be used as a package or cli tool.
It serves as a data loader to support ELT pipelines or any kind of process that requires a heavy data load.

## CLI

This package allows the user to load data from one or multiple json files into a mongodb database.

### Ping database
The `ping` command does a ping in database and returns a connection check.


```bash
mongoloader ping --conn-uri "$MONGO_CONN_URI"
```

### Check if a collection exists
The `collxst` command does a ping in database and returns a connection check.


```bash
mongoloader collxst \
	--conn-uri "$MONGO_CONN_URI" \
	--db-name "$MONGO_DBNAME" \
	--collection "$MONGO_COLLECTION" \
	--app-name "$APPNAME"
```


### Loads in batches - Read and streams file content into a pool of workers (async)
The `load-batch` command iterates over a directory and read/marshal files into documents and insert them into a mongodb database.


```bash
mongoloader load-batch \
		--conn-uri "$MONGO_CONN_URI" \
		--db-name "$MONGO_DBNAME" \
		--collection "$MONGO_COLLECTION" \
		--app-name "$APPNAME" \
		--search-path "./data" \
		--file-prefix "some_prefix" \
		--num-concurrent-files 10
```

### Extract data
The `load` command load a file or set of files into an slice of documents and inserts it into mongodb database.


```bash
mongoloader load \
	--conn-uri "$MONGO_CONN_URI" \
	--db-name "$MONGO_DBNAME" \
	--collection "$MONGO_COLLECTION" \
	--app-name "$APPNAME" \
```
