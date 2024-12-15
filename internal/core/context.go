// trace_context.go
package core

import (
	"bytes"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/null-bd/logger/types"
)

type traceStore struct {
	mu     sync.RWMutex
	traces map[uint64]*traceEntry
}

type traceEntry struct {
	fields    types.Fields
	timestamp time.Time
}

var (
	globalTraceContext *traceStore
	// Configure these based on your needs
	cleanupInterval = 5 * time.Minute
	maxTraceAge     = 30 * time.Minute
	maxTraceEntries = 10000 // Prevent unbounded growth
)

func init() {
	globalTraceContext = &traceStore{
		traces: make(map[uint64]*traceEntry),
	}
	startCleanupRoutine()
}

func startCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			cleanupTraces()
		}
	}()
}

func cleanupTraces() {
	now := time.Now()
	expiredIDs := make([]uint64, 0)

	globalTraceContext.mu.RLock()
	// First pass: identify expired entries
	for gID, entry := range globalTraceContext.traces {
		if now.Sub(entry.timestamp) > maxTraceAge {
			expiredIDs = append(expiredIDs, gID)
		}
	}
	totalEntries := len(globalTraceContext.traces)
	globalTraceContext.mu.RUnlock()

	if len(expiredIDs) > 0 {
		globalTraceContext.mu.Lock()
		// Second pass: remove expired entries
		for _, gID := range expiredIDs {
			delete(globalTraceContext.traces, gID)
		}
		globalTraceContext.mu.Unlock()
	}

	// If still too many entries after age-based cleanup,
	// perform size-based cleanup
	if totalEntries > maxTraceEntries {
		performSizeBasedCleanup()
	}
}

func performSizeBasedCleanup() {
	globalTraceContext.mu.Lock()
	defer globalTraceContext.mu.Unlock()

	// If size is still over limit
	if len(globalTraceContext.traces) > maxTraceEntries {
		// Create slice of entries for sorting
		entries := make([]struct {
			gID       uint64
			timestamp time.Time
		}, 0, len(globalTraceContext.traces))

		for gID, entry := range globalTraceContext.traces {
			entries = append(entries, struct {
				gID       uint64
				timestamp time.Time
			}{gID, entry.timestamp})
		}

		// Sort by timestamp, oldest first
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].timestamp.Before(entries[j].timestamp)
		})

		// Remove oldest entries until we're under the limit
		for i := 0; i < len(entries)-maxTraceEntries; i++ {
			delete(globalTraceContext.traces, entries[i].gID)
		}
	}
}

// SetTraceFields sets trace fields for current goroutine
func SetTraceFields(fields types.Fields) {
	gID := getGoroutineID()
	globalTraceContext.mu.Lock()
	globalTraceContext.traces[gID] = &traceEntry{
		fields:    fields,
		timestamp: time.Now(),
	}
	globalTraceContext.mu.Unlock()
}

// GetTraceFields gets trace fields for current goroutine
func GetTraceFields() types.Fields {
	gID := getGoroutineID()
	globalTraceContext.mu.RLock()
	defer globalTraceContext.mu.RUnlock()

	if entry, exists := globalTraceContext.traces[gID]; exists {
		return entry.fields
	}
	return nil
}

// ClearTraceFields removes trace fields for current goroutine
func ClearTraceFields() {
	gID := getGoroutineID()
	globalTraceContext.mu.Lock()
	delete(globalTraceContext.traces, gID)
	globalTraceContext.mu.Unlock()
}

// getGoroutineID implementation remains the same
func getGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
