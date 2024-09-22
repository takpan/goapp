package httpsrv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncStats(t *testing.T) {
	// Create server
	strChan := make(chan string, 100)
	server := New(strChan)

	// Suppose 2 watchers
	watcherId1 := "watcher_1"
	watcherId2 := "watcher_2"

	// Increment stats for both watchers
	for i := 0; i < 5; i++ {
		server.incStats(watcherId1)
		server.incStats(watcherId2)
	}

	// Increment stats for watcher_2 an additional 3 times
	for i := 0; i < 3; i++ {
		server.incStats(watcherId2)
	}

	// Check if sessionStats has at least 2 entries
	if len(server.sessionStats) < 2 {
		t.Fatalf("Expected at least 2 sessionStats entries, got %d", len(server.sessionStats))
	}

	// Assert that the sent counts are correct
	assert.Equal(t, server.sessionStats[0].sent, 5, "Expected 5 'sent' for watcherId1, got %d", server.sessionStats[0].sent)
	assert.Equal(t, server.sessionStats[1].sent, 8, "Expected 8 'sent' for watcherId1, got %d", server.sessionStats[1].sent)
}
