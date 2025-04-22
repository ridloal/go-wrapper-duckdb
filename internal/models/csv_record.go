package models

type CSVRecord struct {
    ID        int    `json:"id"`
    Filename  string `json:"filename"`
    TableName string `json:"table_name"`
    Status    string `json:"status"`
    CreatedAt string `json:"created_at"`
}