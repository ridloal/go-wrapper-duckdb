package api

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ridloal/go-wrapper-duckdb/internal/database"
	"github.com/ridloal/go-wrapper-duckdb/internal/state"
)

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

type ImportCSVRequest struct {
	DestinationType  string `json:"destination_type"`                  // "duckdb" atau "postgres"
	ConnectionString string `json:"connection_string,omitempty"`       // untuk postgres
	Schema           string `json:"schema,omitempty" default:"public"` // untuk postgres
}

func ImportCSV(db *database.DuckDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if import is already running
		if state.IsImportRunning() {
			return c.Status(409).JSON(fiber.Map{
				"error":  "Import process is already running",
				"status": state.GetImportStatus(),
			})
		}

		var req ImportCSVRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		files, err := filepath.Glob("csv_shared/*.csv")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Error mencari file CSV: %v", err),
			})
		}

		if len(files) == 0 {
			return c.Status(404).JSON(fiber.Map{
				"error": "Tidak ada file CSV di folder csv_shared/",
			})
		}

		// Initialize import status
		state.SetImportStatus(state.ImportStatus{
			IsRunning:      true,
			StartTime:      time.Now(),
			TotalFiles:     len(files),
			ProcessedFiles: 0,
		})

		// Start async import process with destination config
		go processCSVFiles(db, files, req)

		return c.JSON(fiber.Map{
			"message": "Import process started",
			"status":  state.GetImportStatus(),
		})
	}
}

func processCSVFiles(db *database.DuckDB, files []string, req ImportCSVRequest) {
	logger := log.New(os.Stdout, "[CSV Import] ", log.LstdFlags)

	defer func() {
		status := state.GetImportStatus()
		status.IsRunning = false
		state.SetImportStatus(status)
	}()

	var queries []string

	// Setup berdasarkan tipe destinasi
	if req.DestinationType == "postgres" {
		// Setup untuk PostgreSQL
		queries = append(queries, "INSTALL postgres;")
		queries = append(queries, "LOAD postgres;")
		
		// Tambahkan query create schema jika belum ada
		if req.Schema != "" {
			queries = append(queries, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", req.Schema))
		}
		
		queries = append(queries, fmt.Sprintf("ATTACH '%s' AS postgres_db (TYPE postgres, SCHEMA '%s');", 
			req.ConnectionString, 
			req.Schema))
	}

	// Proses setiap file CSV
	for _, file := range files {
		filename := filepath.Base(file)
		tableName := strings.TrimSuffix(filename, ".csv")
		tableName = strings.ReplaceAll(tableName, " ", "_")

		var query string
		if req.DestinationType == "postgres" {
			// Import ke PostgreSQL
			query = fmt.Sprintf(
				"CREATE TABLE postgres_db.%s AS SELECT * FROM read_csv_auto('%s', ignore_errors=true);",
				tableName,
				filepath.ToSlash(file),
			)
		} else {
			// Import ke DuckDB (default)
			query = fmt.Sprintf(
				"CREATE TABLE %s AS SELECT * FROM read_csv_auto('%s', ignore_errors=true);",
				tableName,
				filepath.ToSlash(file),
			)
		}
		queries = append(queries, query)
	}

	// Jalankan semua query dalam satu eksekusi
	combinedQuery := strings.Join(queries, " ")
	if err := db.ExecuteQuery(combinedQuery); err != nil {
		logger.Printf("Error executing combined query: %v", err)
		status := state.GetImportStatus()
		status.LastError = err.Error()
		state.SetImportStatus(status)
		return
	}

	logger.Println("Import process completed successfully")
}

func GetImportStatus(c *fiber.Ctx) error {
	return c.JSON(state.GetImportStatus())
}
