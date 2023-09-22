package cmd

import (
	"os"

	constants "github.com/farovictor/MongoPipeline/src/constants"
	logger "github.com/farovictor/MongoPipeline/src/logging"
	"github.com/spf13/cobra"
)

var mapping string
var connUri string
var dbName string
var appName string
var outputFilePrefix string
var outputPath string
var batchSize int32
var query string
var collectionName string
var numConcurrentFiles int32

// Root Command (does nothing, only prints nice things)
var rootCmd = &cobra.Command{
	Short: "This project aims to support mongodb extractors/loaders",
	Run: func(cmd *cobra.Command, args []string) {
		println("s2")
	},
}

// Extract Command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "This is a extractor for mongodb routines",
	Run:   extractMapping,
}

// Extract Batch Command
var extractBatchesCmd = &cobra.Command{
	Use:   "extract-batch",
	Short: "This is a batch extractor for mongodb routines",
	Run:   extractBatches,
}

// Ping Command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "This is a ping check for mongodb connection",
	Run:   pingExecute,
}

// Check if a collection exists
var collExistsCmd = &cobra.Command{
	Use:   "collxst",
	Short: "This command checks if a defined collection exists",
	Run:   collExistsExecute,
}

// Executes cli
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.ErrorLogger.Printf("%v %s", os.Stderr, err)
		println()
		os.Exit(1)
	}
}

func init() {

	// Root command flags setup
	rootCmd.PersistentFlags().StringVarP(&connUri, "conn-uri", "c", "", "Connection uri for mongodb")
	rootCmd.PersistentFlags().StringVarP(&dbName, "db-name", "d", "", "Database name")
	rootCmd.PersistentFlags().StringVarP(&appName, "app-name", "a", "", "App name")
	rootCmd.MarkFlagsRequiredTogether("conn-uri", "db-name", "app-name")
	// Extract command flags setup
	extractCmd.PersistentFlags().StringVarP(&mapping, "mapping", "m", "", "Mapping name to use for extraction")
	extractCmd.PersistentFlags().StringVarP(&outputFilePrefix, "output-prefix", "o", constants.MappingDefault, "Output filename prefix")
	extractCmd.PersistentFlags().StringVarP(&outputPath, "output-path", "p", ".", "Output folder path")
	extractCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "WHERE clause to attach to query in a valid mongodb syntax")
	extractCmd.PersistentFlags().Int32VarP(&batchSize, "chunk-size", "s", 100, "Chunk size for exported files")
	extractCmd.PersistentFlags().StringVar(&collectionName, "collection", "", "Specify the collection you want to check")
	extractCmd.MarkFlagRequired("collection")
	extractCmd.MarkFlagRequired("mapping")
	// Extract Batches
	extractBatchesCmd.PersistentFlags().StringVarP(&mapping, "mapping", "m", "", "Mapping name to use for extraction")
	extractBatchesCmd.PersistentFlags().StringVarP(&outputFilePrefix, "output-prefix", "o", constants.MappingDefault, "Output filename prefix")
	extractBatchesCmd.PersistentFlags().StringVarP(&outputPath, "output-path", "p", ".", "Output folder path")
	extractBatchesCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "WHERE clause to attach to query in a valid mongodb syntax")
	extractBatchesCmd.PersistentFlags().Int32VarP(&batchSize, "chunk-size", "s", 100, "Chunk size for exported files")
	extractBatchesCmd.PersistentFlags().StringVar(&collectionName, "collection", "", "Specify the collection you want to check")
	extractBatchesCmd.PersistentFlags().Int32VarP(&numConcurrentFiles, "num-concurrent-files", "n", 50, "Number of concurrent files to dump")
	extractBatchesCmd.MarkFlagRequired("collection")
	extractBatchesCmd.MarkFlagRequired("mapping")
	// Collection exists command flags setup
	collExistsCmd.PersistentFlags().StringVar(&collectionName, "collection", "", "Specify the collection you want to check")
	collExistsCmd.MarkFlagRequired("collection")

	// Attaching commands to root
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(extractBatchesCmd)
	rootCmd.AddCommand(pingCmd)
	rootCmd.AddCommand(collExistsCmd)
}
