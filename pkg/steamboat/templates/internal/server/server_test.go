package server

import (
	"testing"
)

func TestNewWithOptions(t *testing.T) {
	srv := NewWithOptions(false)
	if srv == nil {
		t.Error("NewWithOptions() returned nil")
	}

	if srv.db == nil {
		t.Error("database service not initialized")
	}

	if srv.handlers == nil {
		t.Error("handlers not initialized")
	}

	if srv.server == nil {
		t.Error("http server not initialized")
	}
}