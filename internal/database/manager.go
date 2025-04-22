package database

import (
	"fmt"
	"os"
	"path/filepath"
)

const DatabaseDir = "database_duckdb"

// EnsureDatabaseDir memastikan direktori database ada
func EnsureDatabaseDir() error {
	if err := os.MkdirAll(DatabaseDir, 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori database: %v", err)
	}
	return nil
}

// ListDatabases mengembalikan daftar file database
func ListDatabases() ([]string, error) {
	pattern := filepath.Join(DatabaseDir, "*.duckdb")
	return filepath.Glob(pattern)
}