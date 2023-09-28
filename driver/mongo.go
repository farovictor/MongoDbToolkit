package mongo

import (
	"context"
	"sync"
	"time"

	logger "github.com/farovictor/MongodbDriver/logging"
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

func (m *ConnectionHandler) ExtractResults(mapping string, filePrefix string, fileLocation string, process func([]*bson.M, string, string, string) error, coll *mongo.Collection, filter interface{}, opts ...*options.FindOptions) error {

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

	if err := process(results, mapping, filePrefix, fileLocation); err != nil {
		return err
	}
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

// Retrieve collection
func (m *ConnectionHandler) GetCollection(collectionName string) *mongo.Collection {
	logger.InfoLogger.Println("GetCollectionName:", collectionName)
	return m.database.Collection(collectionName)
}

// Streaming results into a pool of workers
// Process function should loop through channel. Once the channel is closed, worker should send a Done signal when finished iterating the channel.
func (m *ConnectionHandler) StreamingResults(mapping, filePrefix, fileLocation string,
	batchSize, numWorkers int32,
	process func(ctx context.Context, batchData <-chan []*bson.M, mapping string, wg *sync.WaitGroup, filePrefix string, folder string),
	coll *mongo.Collection, filter interface{}, opts ...*options.FindOptions) error {

	ctx := context.Background()

	// Creating channel that will handler the results list
	pipe := make(chan []*bson.M, numWorkers)

	// Wait Group to Handler file dumps
	var wg sync.WaitGroup

	// Creating workers and attaching channel
	for i := int32(0); i < numWorkers; i++ {
		wg.Add(1)

		// Dispatching function to workers
		go process(ctx, pipe, mapping, &wg, filePrefix, fileLocation)
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

// This method reads files, populates a slice and inserts into a collection (this may require a huge amount of memory)
func (m *ConnectionHandler) InsertFromFiles(filePrefix, folder string, walk func(filePrefix string, folder string) ([]any, error), coll *mongo.Collection, opts ...*options.InsertManyOptions) error {
	ctx := context.Background()

	// Walking function
	documents, err := walk(filePrefix, folder)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	if len(documents) == 0 {
		logger.WarningLogger.Println("No documents to insert.")
		return nil
	}

	// Inserting resutls
	results, err := coll.InsertMany(ctx, documents, opts...)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	logger.InfoLogger.Printf("Inserted %d documents", len(results.InsertedIDs))
	return nil
}

// This method reads files concurrently, send them to workers that will inserts the values concurrently
func (m *ConnectionHandler) ConcurrentBatchInsert(
	filePrefix, folder string,
	numWorkers int32,
	insert func(ctx context.Context, files <-chan string, wg *sync.WaitGroup, coll *mongo.Collection),
	coll *mongo.Collection) error {

	ctx := context.Background()

	// Creating channel that will handler the results list
	pipe := make(chan string, numWorkers)

	// Wait Group to Handler file dumps
	var wg sync.WaitGroup

	// Creating workers and attaching channel
	for i := int32(0); i < numWorkers; i++ {
		wg.Add(1)
		// dispatching workers
		go insert(ctx, pipe, &wg, coll)
	}

	// Send files across channel
	counter, err := findFiles(filePrefix, folder, pipe)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	// Close files channel
	close(pipe)

	// Wait for all workers
	wg.Wait()

	logger.InfoLogger.Printf("%d files were collected", counter)

	return nil
}
