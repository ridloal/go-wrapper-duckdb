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

func ImportCSV(db *database.DuckDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if import is already running
		if state.IsImportRunning() {
			return c.Status(409).JSON(fiber.Map{
				"error":  "Import process is already running",
				"status": state.GetImportStatus(),
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

		// Start async import process
		go processCSVFiles(db, files)

		return c.JSON(fiber.Map{
			"message": "Import process started",
			"status":  state.GetImportStatus(),
		})
	}
}

func processCSVFiles(db *database.DuckDB, files []string) {
	logger := log.New(os.Stdout, "[CSV Import] ", log.LstdFlags)

	defer func() {
		status := state.GetImportStatus()
		status.IsRunning = false
		state.SetImportStatus(status)
	}()

	for i, file := range files {
		filename := filepath.Base(file)
		logger.Printf("Processing file %d/%d: %s", i+1, len(files), filename)

		tableName := strings.TrimSuffix(filename, ".csv")
		tableName = strings.ReplaceAll(tableName, " ", "_")

		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s AS SELECT * FROM read_csv_auto('%s', ignore_errors=true);",
			tableName, filepath.ToSlash(file))

		if err := db.ExecuteQuery(query); err != nil {
			logger.Printf("Error importing %s: %v", filename, err)
			status := state.GetImportStatus()
			status.LastError = err.Error()
			state.SetImportStatus(status)
			continue
		}

		logger.Printf("Successfully imported %s to table %s", filename, tableName)

		status := state.GetImportStatus()
		status.ProcessedFiles++
		state.SetImportStatus(status)
	}

	logger.Println("Import process completed")
}

func GetImportStatus(c *fiber.Ctx) error {
	return c.JSON(state.GetImportStatus())
}
