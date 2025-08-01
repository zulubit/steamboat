package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"<<!.ProjectName!>>/internal/database/models"
)

type mockDB struct{}

func (m *mockDB) Close() error                    { return nil }
func (m *mockDB) Example() *models.ExampleQueries { return nil }

func TestHelloWorldHandler(t *testing.T) {
	h := New(&mockDB{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.HelloWorldHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	expected := `{"message":"Hello World"}`
	if w.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, w.Body.String())
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
}