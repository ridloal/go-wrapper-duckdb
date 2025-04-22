package database

import (
	"fmt"
	"os"
	"os/exec"
)

// DuckDB menangani operasi database dasar
type DuckDB struct {
	DBPath string
}

// NewDuckDB membuat instance DuckDB baru
func NewDuckDB(dbPath string) *DuckDB {
	return &DuckDB{DBPath: dbPath}
}

// CreateInteractiveCommand membuat command untuk mode interaktif
func (db *DuckDB) CreateInteractiveCommand() *exec.Cmd {
	if db.DBPath != "" {
		return exec.Command("duckdb", db.DBPath)
	}
	return exec.Command("duckdb")
}

// ExecuteQuery menjalankan query tunggal
func (db *DuckDB) ExecuteQuery(query string) error {
	var cmd *exec.Cmd
	if db.DBPath != "" {
		cmd = exec.Command("duckdb", db.DBPath, "-c", query)
	} else {
		cmd = exec.Command("duckdb", "-c", query)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ExecuteFile menjalankan query dari file
func (db *DuckDB) ExecuteFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error membaca file: %v", err)
	}

	return db.ExecuteQuery(string(content))
}
