package mongo

import (
	"fmt"
	"log"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func setup(t testing.TB) (*ConnectionHandler, func(), func(tb testing.TB)) {
	connUri := os.Getenv("MONGO_CONN_URI")
	dbname := os.Getenv("MONGO_DBNAME")

	// Connection setup
	handler := ConnectionHandler{}
	disconnect := handler.Connect(connUri, "Tests")
	handler.setDatabase(dbname)

	return &handler, disconnect, func(tb testing.TB) {
		// Placeholder for teardown scripts
		log.Println("Executing Teardown")
	}
}

func TestPing(t *testing.T) {
	connUri := os.Getenv("MONGO_CONN_URI")
	ok, err := Ping(connUri)

	// Something went wrong when pinging server
	if err != nil {
		t.Error(err)
	}

	// It didnt throw any error but it didnt connect too
	if !ok {
		t.Fail()
	}
}

func TestCollectionExist(t *testing.T) {
	handler, disconnect, teardown := setup(t)
	defer disconnect()
	defer teardown(t)

	collection := os.Getenv("MONGO_COLLECTION")
	if !handler.CollectionExists(collection) {
		t.Errorf("Could not find `%s` collection\n", collection)
	}
}

func TestInsertMany(t *testing.T) {
	handler, disconnect, teardown := setup(t)
	defer disconnect()
	defer teardown(t)

	tempDir := t.TempDir()
	filePrefix := "test_insert"

	err := createFiles(filePrefix, tempDir)
	if err != nil {
		t.Fatal("Error while trying to create files", err)
	}

	opts := options.InsertManyOptions{}

	collName := os.Getenv("MONGO_COLLECTION_TARGET")
	collection := handler.GetCollection(collName)
	fmt.Printf("Collection name %s - Collection %v", collName, collection)

	err = handler.InsertFromFiles(filePrefix, tempDir, ReadFiles, collection, &opts)

	if err != nil {
		t.Errorf("Could not insert on collection: `%s`\n%s", collName, err)
	}
}
