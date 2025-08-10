package database

import (
	"testing"
)

func TestEmailDBStruct(t *testing.T) {
	// Skip this test as it requires FTS5 support which may not be available
	t.Skip("Skipping unit test - FTS5 integration test is in db_integration_test.go")
}
