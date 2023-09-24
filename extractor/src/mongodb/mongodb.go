package mongo

import (
	"context"
	"sync"
	"time"

	files "github.com/farovictor/MongoDbExtractor/src/files"
	logger "github.com/farovictor/MongoDbExtractor/src/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Ping(connUri string) (bool, error) {

	// Create a Client to a MongoDB server and use Ping to verify that the
	// server is running.

	clientOpts := options.Client().ApplyURI(connUri).SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		// log.Fatal(err)
		return false, err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			logger.WarningLogger.Println(err)
		}
	}()

	// Call Ping to verify that the deployment is up and the Client was
	// configured successfully. As mentioned in the Ping documentation, this
	// reduces application resiliency as the server may be temporarily
	// unavailable when Ping is called.
	if err = client.Ping(context.TODO(), readpref.Nearest()); err != nil {
		return false, err
	}

	return true, nil

}

func NewConnectionHandler(connUri string, dbName string, appName string) (*ConnectionHandler, func()) {
	handler := ConnectionHandler{}
	disconnect := handler.Connect(connUri, appName)
	handler.setDatabase(dbName)
	return &handler, disconnect
}

// Connection Handler holds a client and database instances
type ConnectionHandler struct {
	// some fields in here
	client   *mongo.Client
	database *mongo.Database
}

func (m *ConnectionHandler) Connect(connUri string, appName string) func() {
	clientOpts := options.Client().ApplyURI(connUri).SetConnectTimeout(5 * time.Second).SetAppName(appName)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}
	logger.InfoLogger.Println("Client connected!")

	m.client = client
	return func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			logger.ErrorLogger.Fatal(err)
		}
		logger.InfoLogger.Println("Client disconnected!")
	}
}

func (m *ConnectionHandler) setDatabase(dbname string) {
	logger.InfoLogger.Println("Setting database", dbname)
	m.database = m.client.Database(dbname)
}

func (m *ConnectionHandler) ExtractResults(mapping string, filePrefix string, fileLocation string, coll *mongo.Collection, filter interface{}, opts ...*options.FindOptions) error {

	// Find all documents in which the "name" field is "Bob".
	// Specify the Sort option to sort the returned documents by age in
	// ascending order.
	// e.g.: opts := options.Find().SetSort(bson.D{{"age", 1}})
	// e.g.: filter := bson.D{{"name", "Bob"}}
	ctx := context.TODO()
	cursor, err := coll.Find(ctx, filter, opts...)
	defer cursor.Close(ctx)
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	// Get a list of all returned documents and print them out.
	// See the mongo.Cursor documentation for more examples of using cursors.
	var results []*bson.M
	if err = cursor.All(ctx, &results); err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	files.DumpToJsonFile(results, mapping, filePrefix, fileLocation)
	return nil
}

func (m *ConnectionHandler) CollectionExists(collection string) bool {

	logger.InfoLogger.Println("Collection Exists called")
	cursor, err := m.database.ListCollections(context.TODO(), bson.D{})
	if err != nil {
		logger.InfoLogger.Printf("Could not list collections. Error: %s", err)
		return false
	}

	var results []bson.M

	for cursor.Next(context.TODO()) {
		var element bson.M
		if err := cursor.Decode(&element); err != nil {
			logger.ErrorLogger.Printf("Could not decode results. Error: %s", err)
			return false
		}
		results = append(results, element)
	}

	for _, item := range results {
		val, ok := item["name"]
		if ok && val == collection {
			return true
		}
	}
	return false
}

func (m *ConnectionHandler) GetCollection(collectionName string) *mongo.Collection {
	logger.InfoLogger.Println("GetCollectionName:", collectionName)
	return m.database.Collection(collectionName)
}

// Streaming results into a pool of workers
// Check readings:
// Medium article: https://teivah.medium.com/a-closer-look-at-go-sync-package-9f4e4a28c35a
// Pkg doc: https://pkg.go.dev/sync
func (m *ConnectionHandler) StreamingResults(mapping string, batchSize int32, numWorkers int32, filePrefix string, fileLocation string, coll *mongo.Collection, filter interface{}, opts ...*options.FindOptions) error {

	ctx := context.TODO()

	// Creating channel that will handler the results list
	pipe := make(chan []*bson.M, numWorkers)

	// Wait Group to Handler file dumps
	var wg sync.WaitGroup

	// Attaching workers to channel
	for i := int32(0); i < numWorkers; i++ {
		wg.Add(1)
		go files.DumpStreams(ctx, pipe, mapping, &wg, filePrefix, fileLocation)
	}

	cursor, err := coll.Find(ctx, filter, opts...)
	defer cursor.Close(ctx)
	if err != nil {
		// logger.ErrorLogger.Fatal(err)
		close(pipe)
		return err
	}

	// Get a list of all returned documents and print them out.
	// See the mongo.Cursor documentation for more examples of using cursors.
	var counter int32 = 0
	var results []*bson.M
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			// logger.WarningLogger.Println("Error decoding document:", err)
			close(pipe)
			return err
		} else {
			// Appending to slice
			results = append(results, &result)
			counter += 1
		}

		if counter > batchSize {

			// Reset counter
			counter = 0

			// Send data to channel
			pipe <- results

			// Resetting results
			results = []*bson.M{}

		}

	}

	// Send residual results to channel
	pipe <- results

	close(pipe)

	// Wait for all workers
	wg.Wait()

	if cursor.Err() != nil {
		// logger.ErrorLogger.Fatalln("Error while iterating cursor:", cursor.Err())
		return err
	}

	return nil
}
