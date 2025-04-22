package state

import (
	"sync"
	"time"
)

type ImportStatus struct {
	IsRunning   bool      `json:"is_running"`
	StartTime   time.Time `json:"start_time"`
	TotalFiles  int      `json:"total_files"`
	ProcessedFiles int    `json:"processed_files"`
	LastError   string   `json:"last_error,omitempty"`
}

var (
	importStatus ImportStatus
	importMutex  sync.Mutex
)

func GetImportStatus() ImportStatus {
	importMutex.Lock()
	defer importMutex.Unlock()
	return importStatus
}

func SetImportStatus(status ImportStatus) {
	importMutex.Lock()
	defer importMutex.Unlock()
	importStatus = status
}

func IsImportRunning() bool {
	importMutex.Lock()
	defer importMutex.Unlock()
	return importStatus.IsRunning
}