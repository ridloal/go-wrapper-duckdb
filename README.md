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

#### e. Mode Server
Menjalankan aplikasi dalam mode server (API):
```bash
duckdb-wrapper.exe -server
```

### 3. Struktur Folder
- `csv_data/`: Tempat menyimpan file CSV yang akan diimport
- `database_duckdb/`: Tempat menyimpan file database DuckDB
- `cmd/`: Berisi file main.go
- `internal/`: Berisi implementasi internal program
- `csv_shared/`: Tempat menyimpan file CSV untuk mode server

### 4. Catatan Penting
- Pastikan DuckDB CLI sudah terinstall dan dapat diakses dari command line
- File CSV yang akan diimport harus ditempatkan di folder `csv_data/`
- File database akan disimpan di folder `database_duckdb/`
- Gunakan semicolon (;) untuk mengakhiri query dalam mode interaktif
- Untuk mode server, file CSV harus ditempatkan di folder `csv_shared/`
- Server API berjalan pada port 3000
- Import CSV melalui API berjalan secara asynchronous
- Nama tabel akan dibuat otomatis dari nama file CSV (spasi diubah menjadi underscore)

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

### 6. Mode Server
Dengan database spesifik:
```bash
duckdb-wrapper.exe -db database_duckdb/mydb.duckdb -server
```

Server akan berjalan pada port 3000 dengan endpoint berikut:

#### Health Check

**Endpoint**: GET /health  
**Deskripsi**: Mengecek status server

```bash
curl -X GET http://localhost:3000/health
```

**Response**:
```json
{
  "status": "healthy",
  "timestamp": "2024-02-20T10:00:00Z"
}
```

#### Import CSV

**Endpoint**: POST /import-csv  
**Deskripsi**: Mengimport semua file CSV dari folder csv_shared/ secara asynchronous  
**Proses**:
- Mengubah spasi dalam nama file menjadi underscore untuk nama tabel
- Mengabaikan error pada saat import dengan parameter ignore_errors=true

```bash
curl -X POST http://localhost:3000/import-csv
```

**Response Success**:
```json
{
  "message": "Import process started",
  "status": {
    "is_running": true,
    "start_time": "2024-02-20T10:00:00Z",
    "total_files": 2,
    "processed_files": 0
  }
}
```

**Response Error** (jika proses import sedang berjalan):
```json
{
  "error": "Import process is already running",
  "status": {
    "is_running": true,
    "start_time": "2024-02-20T10:00:00Z",
    "total_files": 2,
    "processed_files": 1
  }
}
```

#### Import Status

**Endpoint**: GET /import-status  
**Deskripsi**: Mendapatkan status proses import yang sedang berjalan

```bash
curl http://localhost:3000/import-status
```

**Response**:
```json
{
  "is_running": true,
  "start_time": "2024-02-20T10:00:00Z",
  "total_files": 2,
  "processed_files": 1,
  "last_error": "error message if any"
}
```