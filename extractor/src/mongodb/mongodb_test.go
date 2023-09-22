package mongo

import (
	"log"
	"os"
	"testing"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
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

// func TestExtract(t *testing.T) {
// 	handler, teardown := setup(t)
// 	defer teardown(t)

// 	collname := "army"
// 	coll := handler.database.Collection(collname)

// 	// Test
// 	opts := options.Find().SetSort(bson.D{{"name", 1}})
// 	results := handler.Extract(coll, bson.D{{"type", "orc"}}, opts)
// 	if len(*results) < 100 {
// 		t.Error("Less orcs then expected")
// 	}
// }

// func TestSerializeToParquet(t *testing.TB) {
//
// }
