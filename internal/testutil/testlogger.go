package testutil

import (
	"fmt"
	"sync"
)

type LogEntry struct {
	Level   string
	Message string
}

type TestLogger struct {
	mu      sync.Mutex
	entries []LogEntry
}

func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

func (l *TestLogger) Info(format string, args ...interface{}) {
	l.append("info", format, args...)
}

func (l *TestLogger) Error(format string, args ...interface{}) {
	l.append("error", format, args...)
}

func (l *TestLogger) Debug(format string, args ...interface{}) {
	l.append("debug", format, args...)
}

func (l *TestLogger) Entries() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()

	entries := make([]LogEntry, len(l.entries))
	copy(entries, l.entries)
	return entries
}

func (l *TestLogger) LastEntry() (LogEntry, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.entries) == 0 {
		return LogEntry{}, false
	}

	return l.entries[len(l.entries)-1], true
}

func (l *TestLogger) append(level string, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = append(l.entries, LogEntry{
		Level:   level,
		Message: fmt.Sprintf(format, args...),
	})
}
