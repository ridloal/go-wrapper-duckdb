# Dokumentasi Go Wrapper DuckDB ETL

## Cara Menjalankan Program

### 1. Build Program
```bash
go build -o duckdb-wrapper.exe cmd/main.go
```

### 2. Mode Penggunaan

#### a. Mode Query Tunggal
Menjalankan query DuckDB tunggal:
```bash
duckdb-wrapper.exe -query "SELECT * FROM table_name;"
```

Dengan database spesifik:
```bash
duckdb-wrapper.exe -db database_duckdb/mydb.duckdb -query "SELECT * FROM table_name;"
```

#### b. Mode File SQL
Menjalankan query dari file SQL:
```bash
duckdb-wrapper.exe -file path/to/queries.sql
```

Dengan database spesifik:
```bash
duckdb-wrapper.exe -db database_duckdb/mydb.duckdb -file path/to/queries.sql
```

#### c. Mode Interaktif
Menjalankan DuckDB dalam mode interaktif:
```bash
duckdb-wrapper.exe -interactive
```

Dengan database spesifik:
```bash
duckdb-wrapper.exe -db database_duckdb/mydb.duckdb -interactive
```

#### d. Mode Import CSV
Mengimport file CSV ke dalam tabel:
```bash
duckdb-wrapper.exe -import-csv
```

Dengan database spesifik:
```bash
duckdb-wrapper.exe -db database_duckdb/mydb.duckdb -import-csv
```

### 3. Struktur Folder
- `csv_data/`: Tempat menyimpan file CSV yang akan diimport
- `database_duckdb/`: Tempat menyimpan file database DuckDB
- `cmd/`: Berisi file main.go
- `internal/`: Berisi implementasi internal program

### 4. Catatan Penting
- Pastikan DuckDB CLI sudah terinstall dan dapat diakses dari command line
- File CSV yang akan diimport harus ditempatkan di folder `csv_data/`
- File database akan disimpan di folder `database_duckdb/`
- Gunakan semicolon (;) untuk mengakhiri query dalam mode interaktif

### 5. Contoh Penggunaan Lengkap
```bash
# Membuat tabel dan mengisi data
duckdb-wrapper.exe -query "CREATE TABLE users (id INTEGER, name VARCHAR);"

# Import data dari CSV
duckdb-wrapper.exe -import-csv

# Query data dalam mode interaktif
duckdb-wrapper.exe -interactive

# Menjalankan file SQL
duckdb-wrapper.exe -file queries/create_tables.sql
```

        