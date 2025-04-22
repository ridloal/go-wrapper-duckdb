package cli

import "flag"

type Flags struct {
    Query       string
    File        string
    Interactive bool
    DBFile      string
    ImportCSV   bool
    Server      bool
}

func ParseFlags() *Flags {
    f := &Flags{}
    flag.StringVar(&f.Query, "query", "", "DuckDB query to execute")
    flag.StringVar(&f.File, "file", "", "SQL file containing DuckDB queries")
    flag.BoolVar(&f.Interactive, "interactive", false, "Run in interactive mode")
    flag.StringVar(&f.DBFile, "db", "", "Path to DuckDB database file")
    flag.BoolVar(&f.ImportCSV, "import-csv", false, "Import CSV file to table")
    flag.BoolVar(&f.Server, "server", false, "Run as HTTP server")
    flag.Parse()
    return f
}
