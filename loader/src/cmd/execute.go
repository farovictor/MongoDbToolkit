package cmd

import (
	"context"
	"os"
	"sync"

	file "github.com/farovictor/MongoDbLoader/src/fs"
	logger "github.com/farovictor/MongoDbLoader/src/logging"
	mongo "github.com/farovictor/MongodbDriver"
	"github.com/spf13/cobra"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

// Execution logic for ping command
func PingExecute(cmd *cobra.Command, args []string) {
	logger.Initialize(logLevel)

	if cmd.Flags().Lookup("conn-uri").Changed {
		ping, err := mongo.Ping(connUri)

		if err != nil {
			logger.ErrorLogger.Println("Error while pinging server:", err)
		}

		if ping {
			logger.InfoLogger.Println("Ping was successful")
		} else {
			logger.WarningLogger.Println("Ping wasn't successful. Check your connection string or network.")
		}
	}
}

// Execution logic for load command
func LoadFile(cmd *cobra.Command, args []string) {
	logger.Initialize(logLevel)

	if collectionName == "" {
		logger.ErrorLogger.Fatalln("No collection specified")
	}

	handler, disconnect := mongo.NewConnectionHandler(connUri, dbName, appName)
	defer disconnect()

	coll := handler.GetCollection(collectionName)
	logger.InfoLogger.Println("Collection retrieved")

	// TODO: Extend this to allow customization
	opts := options.InsertManyOptions{}

	logger.InfoLogger.Println("Processing record")
	if err := handler.InsertFromFiles(filePrefix, searchPath, mongo.ReadFiles, coll, &opts); err != nil {
		logger.ErrorLogger.Fatalln(err)
	}
}

// Execution logic for insert-batches command
func InsertBatches(cmd *cobra.Command, args []string) {

	logger.Initialize(logLevel)

	if collectionName == "" {
		logger.ErrorLogger.Fatalln("No collection specified")
	}

	handler, disconnect := mongo.NewConnectionHandler(connUri, dbName, appName)
	defer disconnect()

	coll := handler.GetCollection(collectionName)

	logger.InfoLogger.Println("Processing files")

	if err := handler.ConcurrentBatchInsert(filePrefix, searchPath, numConcurrentFiles, insertArray, coll); err != nil {
		logger.ErrorLogger.Fatalln(err)
	}
	logger.InfoLogger.Println("Ending Insert Batches")
}

// This function retrieves the files from channel, reads, serialize it and load into a collection
func insertArray(ctx context.Context, files <-chan string, wg *sync.WaitGroup, coll *mongodb.Collection) {
	defer wg.Done()

	for filePath := range files {

		data, err := file.ReadFileToArray(filePath)

		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		results, err := coll.InsertMany(ctx, data, nil)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}

		logger.DebugLogger.Printf("Inserted %d documents\n", len(results.InsertedIDs))
	}
}

// Execution logic for collxst command
func CollExistsExecute(cmd *cobra.Command, args []string) {
	logger.Initialize(logLevel)

	if !cmd.Flags().Lookup("app-name").Changed {
		os.Exit(1)
	}
	if cmd.Flags().Lookup("collection").Changed && cmd.Flags().Lookup("db-name").Changed && cmd.Flags().Lookup("conn-uri").Changed {
		handler, disconnect := mongo.NewConnectionHandler(connUri, dbName, appName)
		defer disconnect()
		logger.InfoLogger.Println(handler)
		exist := handler.CollectionExists(collectionName)
		logger.InfoLogger.Printf("Collection %s exists? %v\n", collectionName, exist)
	}
}
