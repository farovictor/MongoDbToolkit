package mongo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	logger "github.com/farovictor/MongodbDriver/logging"
	"go.mongodb.org/mongo-driver/bson"
)

func createFiles(filePrefix, folder string) error {
	type placeholder struct {
		fileName    string
		fileContent string
	}

	items := []placeholder{
		{"fileA.json", `[{"_id":"65063a2c27b6b3d5da64db70","currency":"SYP","email":"ZijqgSf@bpgmOSH.info","ipv4":"10.150.13.190","ipv6":"328a:559:1249:1612:301d:f2b4:3102:420f","latitude":79.68523406982422,"longitude":144.9792938232422,"month_name":"May","password":"sUxKjlRhHPWBpaxVkgOIUNsKYvADThUVdtwFVJBLNoWXglcTNt","phone":"936-758-1021","url":"http://bULQMDj.org/YcJwcRy","username":"surFtyG","word":"quaerat"}]`},
		{"fileB.json", `[{"_id":"65063a2c27b6b3d5da64db75","currency":"BYN","email":"wpodDGh@rCtHsfo.biz","ipv4":"253.171.52.92","ipv6":"bac2:eae9:525a:b35c:fbf7:9fd7:6617:62e","latitude":58.049644470214844,"longitude":-114.04580688476562,"month_name":"December","password":"DLLOHcFnnDGypmBVtOZJuVZZhcjJmQuPEZBjxvrBDRBDKZKkto","phone":"916-210-4735","url":"https://www.kouwsli.net/","username":"cjpirgG","word":"minima"}]`},
	}

	for _, item := range items {
		filepath := fmt.Sprintf("%s/%s_%s", folder, filePrefix, item.fileName)

		file, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.WriteString(item.fileContent)
		if err != nil {
			return err
		}
	}
	return nil
}

// Function to read files listed in test directory
func ReadFiles(filePrefix string, folder string) ([]any, error) {
	// slice with results
	var documents []any

	// Waling test directory
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}

		// To check if file does not have a regular mode
		if !info.Mode().IsRegular() {
			return nil
		}

		// Read only files that match prefix
		if !info.IsDir() && strings.HasPrefix(info.Name(), filePrefix) {
			file, err := os.Open(path)
			if err != nil {
				logger.ErrorLogger.Println(err)
				return err
			}
			defer file.Close()

			data, err := ioutil.ReadAll(file)

			var jsonArray []map[string]any

			err = json.Unmarshal(data, &jsonArray)
			if err != nil {
				logger.ErrorLogger.Println(err)
				return err
			}

			for _, jsonObj := range jsonArray {
				bsonMap, err := bson.Marshal(jsonObj)
				if err != nil {
					logger.ErrorLogger.Printf("Error converting to BSON: %v\n", err)
					return err
				}
				bsonM := bson.M{}
				err = bson.Unmarshal(bsonMap, &bsonM)
				if err != nil {
					logger.ErrorLogger.Printf("Error unmarshalling BSON: %v\n", err)
					return err
				}
				documents = append(documents, &bsonM)
			}

		}
		return nil
	})

	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return documents, nil
}

func findFiles(filePrefix string, folder string, pipe chan<- string) (int, error) {

	fileCounter := 0

	// Walking inside folder
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// To check if file does not have a regular mode
		if !info.Mode().IsRegular() {
			return nil
		}

		// Read only files that match prefix
		if !info.IsDir() && strings.HasPrefix(info.Name(), filePrefix) {
			pipe <- path
			fileCounter++
		}
		return nil
	})

	if err != nil {
		return fileCounter, err
	}

	return fileCounter, nil
}
