package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	
	"<<!.ProjectName!>>/internal/database"
)

func TestHomeHandler(t *testing.T) {
	// Use a temporary in-memory database for testing
	db := database.New()
	defer db.Close()
	
	h := New(db)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.HomeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	
	// Check if the component rendered with expected content
	if !strings.Contains(body, "Welcome to Steamboat") {
		t.Errorf("expected body to contain 'Welcome to Steamboat', got %s", body)
	}
	
	if !strings.Contains(body, "<button>Button</button>") {
		t.Errorf("expected body to contain button element, got %s", body)
	}
}