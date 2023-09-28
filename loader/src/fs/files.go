package fs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	logger "github.com/farovictor/MongoDbLoader/src/logging"
	"go.mongodb.org/mongo-driver/bson"
)

// Reads file and returns a bson.M array
func ReadFileToArray(filePath string) ([]any, error) {

	file, err := os.Open(filePath)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)

	var jsonArray []map[string]any

	err = json.Unmarshal(data, &jsonArray)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	// This is the same as: []*bson.M
	documents := make([]any, len(jsonArray))

	for i, jsonObj := range jsonArray {
		bsonMap, err := bson.Marshal(jsonObj)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		bsonM := bson.M{}
		err = bson.Unmarshal(bsonMap, &bsonM)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		documents[i] = &bsonM
	}
	return documents, nil
}

// Emits filtered files to a channel
func EmitFilesToChannel(filePrefix string, searchPath string, emit chan<- string) error {

	// Walking through directory
	return filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		// To check if file does not have a regular mode
		if !info.Mode().IsRegular() {
			return nil
		}

		// Emit only files that match prefix
		if !info.IsDir() && strings.HasPrefix(info.Name(), filePrefix) {
			emit <- info.Name()
		}
		return nil
	})
}
