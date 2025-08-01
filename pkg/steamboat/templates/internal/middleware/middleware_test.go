package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"<<!.ProjectName!>>/internal/middleware/session"
)

func TestLogger(t *testing.T) {
	handler := Logger()
	if handler == nil {
		t.Error("Logger() returned nil")
	}
}

func TestCORS(t *testing.T) {
	handler := CORS()
	if handler == nil {
		t.Error("CORS() returned nil")
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	corsHandler := handler(testHandler)

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")

	w := httptest.NewRecorder()
	corsHandler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected CORS headers to be set")
	}
}

func TestSessionMiddleware(t *testing.T) {
	// Test with default config
	middleware := Session(nil)
	if middleware == nil {
		t.Fatal("Session middleware returned nil")
	}
	
	// Create a test handler that checks session
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.GetSession(r)
		if sess == nil {
			t.Error("Session should be available in context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		// Session should be new/empty
		if sess.IsAuthenticated() {
			t.Error("New session should not be authenticated")
		}
		
		w.WriteHeader(http.StatusOK)
	})
	
	// Wrap with session middleware
	handler := middleware(testHandler)
	
	// Test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestSessionMiddlewareWithExistingSession(t *testing.T) {
	config := &session.Config{
		CookieName: "test_session",
		SecretKey:  "test-secret-key-for-testing-32ch",
		MaxAge:     3600,
		HttpOnly:   true,
		Secure:     false,
		SameSite:   http.SameSiteLaxMode,
		Path:       "/",
	}
	
	// First, create a session cookie
	originalSession := &session.Session{
		UserID:   123,
		Username: "testuser",
	}
	
	w := httptest.NewRecorder()
	err := session.SaveSession(w, originalSession, config)
	if err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}
	
	// Get the session cookie
	var sessionCookie *http.Cookie
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == config.CookieName {
			sessionCookie = cookie
			break
		}
	}
	
	if sessionCookie == nil {
		t.Fatal("Session cookie not found")
	}
	
	// Now test middleware with existing session
	middleware := Session(config)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.GetSession(r)
		if sess == nil {
			t.Error("Session should be available in context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		if sess.UserID != 123 {
			t.Errorf("Expected UserID 123, got %d", sess.UserID)
		}
		
		if sess.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got %s", sess.Username)
		}
		
		if !sess.IsAuthenticated() {
			t.Error("Session should be authenticated")
		}
		
		w.WriteHeader(http.StatusOK)
	})
	
	handler := middleware(testHandler)
	
	// Create request with session cookie
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(sessionCookie)
	w = httptest.NewRecorder()
	
	handler.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRateLimiter(t *testing.T) {
	handler := RateLimiter(100)
	if handler == nil {
		t.Error("RateLimiter() returned nil")
	}
}

func TestCompress(t *testing.T) {
	handler := Compress()
	if handler == nil {
		t.Error("Compress() returned nil")
	}
}

func TestRequestID(t *testing.T) {
	handler := RequestID()
	if handler == nil {
		t.Error("RequestID() returned nil")
	}
}

func TestRecoverer(t *testing.T) {
	handler := Recoverer()
	if handler == nil {
		t.Error("Recoverer() returned nil")
	}
}