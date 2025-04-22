package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var scanner *bufio.Scanner

func init() {
	scanner = bufio.NewScanner(os.Stdin)
}

// GetUserInput membaca input dari pengguna dengan prompt
func GetUserInput(prompt string, options []string) string {
	fmt.Print(prompt)
	if !scanner.Scan() {
		return ""
	}
	input := scanner.Text()

	if options != nil {
		if idx, err := strconv.Atoi(input); err == nil {
			if idx > 0 && idx <= len(options) {
				return options[idx-1]
			}
		}
		return ""
	}

	return input
}

// SelectDatabaseFile memungkinkan pengguna memilih atau membuat file database
func SelectDatabaseFile() string {
	fmt.Println("Pilih opsi database:")
	fmt.Println("1. Buat database baru")
	fmt.Println("2. Buka database yang ada")
	fmt.Println("3. Gunakan database in-memory (non-persistent)")

	choice := GetUserInput("Masukkan pilihan (1-3): ", nil)

	switch choice {
	case "1":
		filename := GetUserInput("Masukkan nama untuk file database baru: ", nil)
		if !strings.HasSuffix(filename, ".duckdb") {
			filename += ".duckdb"
		}
		return filepath.Join("database_duckdb", filename)

	case "2":
		files, err := filepath.Glob("database_duckdb/*.duckdb")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error mencari file database: %v\n", err)
			return ""
		}

		if len(files) == 0 {
			fmt.Println("Tidak ada file .duckdb di direktori database.")
			path := GetUserInput("Masukkan path ke file database: ", nil)
			return path
		}

		fmt.Println("\nFile database yang tersedia:")
		for i, file := range files {
			fmt.Printf("%d. %s\n", i+1, filepath.Base(file))
		}

		return GetUserInput("Masukkan nomor atau path lengkap: ", files)

	case "3":
		return ""

	default:
		fmt.Println("Pilihan tidak valid, menggunakan database in-memory.")
		return ""
	}
}