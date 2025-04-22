package api

import (
    "fmt"
    "path/filepath"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/ridloal/go-wrapper-duckdb/internal/database"
    "github.com/ridloal/go-wrapper-duckdb/internal/models"
)

func HealthCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

func ImportCSV(db *database.DuckDB) fiber.Handler {
    return func(c *fiber.Ctx) error {
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

        var results []models.CSVRecord
        for _, file := range files {
            filename := filepath.Base(file)
            tableName := filename[:len(filename)-4] // menghapus ekstensi .csv

            // Import CSV ke DuckDB
            query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s AS SELECT * FROM read_csv_auto('%s');",
                tableName, filepath.ToSlash(file))

            err := db.ExecuteQuery(query)
            status := "success"
            if err != nil {
                status = "failed"
            }

            record := models.CSVRecord{
                Filename:  filename,
                TableName: tableName,
                Status:    status,
                CreatedAt: time.Now().Format(time.RFC3339),
            }
            results = append(results, record)
        }

        return c.JSON(results)
    }
}