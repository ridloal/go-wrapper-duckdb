package cli

import (
	"fmt"
	"path/filepath"

	"github.com/ridloal/go-wrapper-duckdb/internal/database"
	"github.com/ridloal/go-wrapper-duckdb/internal/utils"
)

// ImportCSV mengimport file CSV ke dalam tabel DuckDB
func ImportCSV(db *database.DuckDB) error {
	if db.DBPath == "" {
		selectedDB := utils.SelectDatabaseFile()
		db.DBPath = selectedDB
	}

	// Baca file CSV dari folder csv_data
	files, err := filepath.Glob("csv_data/*.csv")
	if err != nil {
		return fmt.Errorf("error mencari file CSV: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("tidak ada file CSV di folder csv_data/")
	}

	fmt.Println("\nFile CSV yang tersedia:")
	for i, file := range files {
		fmt.Printf("%d. %s\n", i+1, filepath.Base(file))
	}

	selectedCsv := utils.GetUserInput("\nPilih nomor file CSV: ", files)
	if selectedCsv == "" {
		return fmt.Errorf("pilihan file CSV tidak valid")
	}

	tableName := utils.GetUserInput("Masukkan nama tabel untuk import: ", nil)
	if tableName == "" {
		return fmt.Errorf("nama tabel tidak boleh kosong")
	}

	// Buat dan eksekusi query import
	query := fmt.Sprintf("CREATE TABLE %s AS SELECT * FROM read_csv_auto('%s');",
		tableName, filepath.ToSlash(selectedCsv))

	if err := db.ExecuteQuery(query); err != nil {
		return fmt.Errorf("error mengimport CSV: %v", err)
	}

	fmt.Printf("\nBerhasil mengimport %s ke tabel %s\n",
		filepath.Base(selectedCsv), tableName)
	return nil
}