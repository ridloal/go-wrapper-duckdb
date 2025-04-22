package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ridloal/go-wrapper-duckdb/internal/cli"
	"github.com/ridloal/go-wrapper-duckdb/internal/database"
)

func main() {
	// Parse command-line flags
	flags := cli.ParseFlags()

	// Validate input
	if flags.Query == "" && flags.File == "" && !flags.Interactive && !flags.ImportCSV {
		fmt.Println("Error: You must specify either -query, -file, -interactive, or -import-csv")
		flag.Usage()
		os.Exit(1)
	}

	// Ensure database directory exists
	if err := database.EnsureDatabaseDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create database instance
	db := database.NewDuckDB(flags.DBFile)

	// Execute based on input mode
	var err error
	switch {
	case flags.ImportCSV:
		err = cli.ImportCSV(db)
	case flags.Query != "":
		err = db.ExecuteQuery(flags.Query)
	case flags.File != "":
		err = db.ExecuteFile(flags.File)
	case flags.Interactive:
		err = cli.RunInteractiveMode(db)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
