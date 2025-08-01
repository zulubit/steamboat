package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"<<!.ProjectName!>>/internal/database"
	"<<!.ProjectName!>>/internal/handlers"
)

func TestSetup(t *testing.T) {
	db := database.New()
	defer db.Close()
	
	h := handlers.New(db)
	router := Setup(h)

	testCases := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/", http.StatusOK},
		{"POST", "/nonexistent", http.StatusNotFound},
		{"GET", "/nonexistent", http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tc.status {
				t.Errorf("expected status %d, got %d", tc.status, w.Code)
			}
		})
	}
}

func TestCORSHeaders(t *testing.T) {
	db := database.New()
	defer db.Close()
	
	h := handlers.New(db)
	router := Setup(h)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected CORS headers to be set")
	}
}