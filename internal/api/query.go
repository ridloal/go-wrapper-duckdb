package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ridloal/go-wrapper-duckdb/internal/database"
)

type QueryRequest struct {
	Query string `json:"query"`
}

func ExecuteQuery(db *database.DuckDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req QueryRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		if req.Query == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Query cannot be empty",
			})
		}

		err := db.ExecuteQuery(req.Query)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Query executed successfully",
		})
	}
}