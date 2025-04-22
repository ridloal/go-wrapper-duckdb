package main

import (
	"flag"
	"fmt"
	"os"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ridloal/go-wrapper-duckdb/internal/api"
	"github.com/ridloal/go-wrapper-duckdb/internal/cli"
	"github.com/ridloal/go-wrapper-duckdb/internal/database"
)

func main() {
	// Parse command-line flags
	flags := cli.ParseFlags()

	// Validate input
	if flags.Query == "" && flags.File == "" && !flags.Interactive && !flags.ImportCSV && !flags.Server {
		fmt.Println("Error: You must specify either -query, -file, -interactive, -import-csv, or -server")
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
	case flags.Server:
		app := fiber.New()
		
		// API routes
		app.Get("/health", api.HealthCheck)
		app.Post("/import-csv", api.ImportCSV(db))
		app.Get("/import-status", api.GetImportStatus)

		log.Fatal(app.Listen(":3000"))
		return
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
