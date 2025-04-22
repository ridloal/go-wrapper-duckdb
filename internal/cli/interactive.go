package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ridloal/go-wrapper-duckdb/internal/database"
	"github.com/ridloal/go-wrapper-duckdb/internal/utils"
)

// RunInteractiveMode menjalankan DuckDB dalam mode interaktif
func RunInteractiveMode(db *database.DuckDB) error {
	var selectedDB string

	if db.DBPath == "" {
		// Jika tidak ada DB yang ditentukan, biarkan pengguna memilih atau membuat
		selectedDB = utils.SelectDatabaseFile()
		db.DBPath = selectedDB
	} else {
		selectedDB = db.DBPath
	}

	fmt.Printf("DuckDB Interactive Mode menggunakan database: %s (ketik 'exit' untuk keluar)\n", selectedDB)
	scanner := bufio.NewScanner(os.Stdin)

	// Gunakan proses DuckDB yang persisten
	cmd := db.CreateInteractiveCommand()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error membuat stdin pipe: %v", err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error memulai DuckDB: %v", err)
	}

	// Buffer untuk query multi-baris
	var queryBuffer strings.Builder

	for {
		if queryBuffer.Len() == 0 {
			fmt.Print("duckdb> ")
		} else {
			fmt.Print("....   > ")
		}

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		if strings.ToLower(input) == "exit" {
			fmt.Fprintln(stdin, ".exit")
			break
		}

		queryBuffer.WriteString(input)

		if strings.HasSuffix(strings.TrimSpace(input), ";") {
			fmt.Fprintln(stdin, queryBuffer.String())
			queryBuffer.Reset()
		} else {
			queryBuffer.WriteString(" ")
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error membaca input: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error: DuckDB keluar dengan: %v", err)
	}

	return nil
}
