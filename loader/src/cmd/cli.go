package cmd

import (
	"fmt"
	"os"

	logger "github.com/farovictor/MongoDbLoader/src/logging"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
	BuildTime string

	connUri            string
	dbName             string
	appName            string
	filePrefix         string
	searchPath         string
	collectionName     string
	numConcurrentFiles int32
	logLevel           string
)

// Root Command (does nothing, only prints nice things)
var rootCmd = &cobra.Command{
	Short:   "This project aims to support mongodb loading pipelines",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("For more info, visit: https://github.com/farovictor/MongoDbToolkit\n")
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Built: %s\n", BuildTime)
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Log-Level: %v\n", logLevel)
	},
}

// Load Command
var loadCmd = &cobra.Command{
	Use:     "load",
	Version: rootCmd.Version,
	Short:   "Loads a json file into mongodb collection",
	Run:     LoadFile,
}

// Batch Load Command
var loadBatchesCmd = &cobra.Command{
	Use:     "load-batch",
	Version: rootCmd.Version,
	Short:   "Loads a set of json files into a mongodb collection (concurrently)",
	Run:     InsertBatches,
}

// Ping Command
var pingCmd = &cobra.Command{
	Use:     "ping",
	Version: rootCmd.Version,
	Short:   "Ping a mongodb server",
	Run:     PingExecute,
}

// Check if a collection exists
var collExistsCmd = &cobra.Command{
	Use:     "collxst",
	Version: rootCmd.Version,
	Short:   "This command checks if a defined collection exists",
	Run:     CollExistsExecute,
}

// Executes cli
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.ErrorLogger.Printf("%v %s\n", os.Stderr, err)
		println()
		os.Exit(1)
	}
}

func init() {

	// Root command flags setup
	rootCmd.PersistentFlags().StringVarP(&connUri, "conn-uri", "c", "", "Connection uri for mongodb")
	rootCmd.PersistentFlags().StringVarP(&dbName, "db-name", "d", "", "Database name")
	rootCmd.PersistentFlags().StringVarP(&appName, "app-name", "a", "", "App name")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Set a max log level")
	rootCmd.MarkFlagsRequiredTogether("conn-uri", "db-name", "app-name")
	// Load command flags setup
	loadCmd.PersistentFlags().StringVarP(&filePrefix, "file-prefix", "o", "", "Filename prefix")
	loadCmd.PersistentFlags().StringVarP(&searchPath, "search-path", "p", ".", "Search path to look for files")
	loadCmd.PersistentFlags().StringVar(&collectionName, "collection", "", "Specify the collection you want to check")
	loadCmd.MarkFlagRequired("collection")
	// Batch Load command flags setup
	loadBatchesCmd.PersistentFlags().StringVarP(&filePrefix, "file-prefix", "o", "", "Filename prefix")
	loadBatchesCmd.PersistentFlags().StringVarP(&searchPath, "search-path", "p", ".", "Search path to look for files")
	loadBatchesCmd.PersistentFlags().StringVar(&collectionName, "collection", "", "Specify the collection you want to check")
	loadBatchesCmd.PersistentFlags().Int32VarP(&numConcurrentFiles, "num-concurrent-files", "n", 50, "Number of concurrent files to dump")
	loadBatchesCmd.MarkFlagRequired("collection")
	// Collection exists command flags setup
	collExistsCmd.PersistentFlags().StringVar(&collectionName, "collection", "", "Specify the collection you want to check")
	collExistsCmd.MarkFlagRequired("collection")

	// Attaching commands to root
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(loadBatchesCmd)
	rootCmd.AddCommand(pingCmd)
	rootCmd.AddCommand(collExistsCmd)
}
