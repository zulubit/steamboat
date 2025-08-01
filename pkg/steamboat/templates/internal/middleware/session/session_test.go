package session

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSession(t *testing.T) {
	sess := NewSession()
	if sess == nil {
		t.Fatal("NewSession returned nil")
	}
	
	if sess.Data == nil {
		t.Error("NewSession should initialize Data map")
	}
	
	if sess.IsAuthenticated() {
		t.Error("New session should not be authenticated")
	}
}

func TestSessionDataOperations(t *testing.T) {
	sess := NewSession()
	
	// Test Set and Get
	sess.Set("key1", "value1")
	if val, exists := sess.Get("key1"); !exists || val != "value1" {
		t.Errorf("Expected key1=value1, got %s (exists: %v)", val, exists)
	}
	
	// Test non-existent key
	if val, exists := sess.Get("nonexistent"); exists || val != "" {
		t.Errorf("Expected empty value for non-existent key, got %s (exists: %v)", val, exists)
	}
	
	// Test Delete
	sess.Delete("key1")
	if _, exists := sess.Get("key1"); exists {
		t.Error("Key should be deleted")
	}
}

func TestSessionAuthentication(t *testing.T) {
	sess := NewSession()
	
	// Initially not authenticated
	if sess.IsAuthenticated() {
		t.Error("New session should not be authenticated")
	}
	
	// Set user ID
	sess.UserID = 123
	sess.Username = "testuser"
	sess.Email = "test@example.com"
	
	if !sess.IsAuthenticated() {
		t.Error("Session with UserID should be authenticated")
	}
	
	// Test Clear
	sess.Clear()
	if sess.IsAuthenticated() || sess.UserID != 0 || sess.Username != "" {
		t.Error("Clear should reset all session data")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	
	if config.CookieName != "steamboat_session" {
		t.Errorf("Expected cookie name 'steamboat_session', got %s", config.CookieName)
	}
	
	if config.MaxAge != 86400*7 {
		t.Errorf("Expected MaxAge to be 7 days, got %d", config.MaxAge)
	}
	
	if !config.HttpOnly {
		t.Error("Expected HttpOnly to be true")
	}
	
	if config.SecretKey == "" {
		t.Error("SecretKey should not be empty")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	secretKey := "test-secret-key-for-encryption"
	testData := []byte("Hello, World! This is test data.")
	
	// Test encryption
	encrypted, err := encrypt(testData, secretKey)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	
	if len(encrypted) == 0 {
		t.Error("Encrypted data should not be empty")
	}
	
	// Test decryption
	decrypted, err := decrypt(encrypted, secretKey)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
	
	if string(decrypted) != string(testData) {
		t.Errorf("Decrypted data doesn't match original. Expected: %s, Got: %s", testData, decrypted)
	}
}

func TestEncryptDecryptWithWrongKey(t *testing.T) {
	secretKey := "test-secret-key"
	wrongKey := "wrong-secret-key"
	testData := []byte("sensitive data")
	
	encrypted, err := encrypt(testData, secretKey)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	
	// Try to decrypt with wrong key
	_, err = decrypt(encrypted, wrongKey)
	if err == nil {
		t.Error("Decryption with wrong key should fail")
	}
}

func TestSaveAndLoadSession(t *testing.T) {
	config := &Config{
		CookieName: "test_session",
		SecretKey:  "test-secret-key-for-testing-32ch",
		MaxAge:     3600,
		HttpOnly:   true,
		Secure:     false,
		SameSite:   http.SameSiteLaxMode,
		Path:       "/",
	}
	
	// Create a test session
	originalSession := &Session{
		UserID:   123,
		Username: "testuser",
		Email:    "test@example.com",
		Data:     map[string]string{"role": "admin", "theme": "dark"},
	}
	
	// Save session
	w := httptest.NewRecorder()
	err := SaveSession(w, originalSession, config)
	if err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}
	
	// Check cookie was set
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies were set")
	}
	
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == config.CookieName {
			sessionCookie = cookie
			break
		}
	}
	
	if sessionCookie == nil {
		t.Fatal("Session cookie not found")
	}
	
	// Load session
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(sessionCookie)
	
	loadedSession, err := LoadSession(req, config)
	if err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}
	
	// Verify loaded session matches original
	if loadedSession.UserID != originalSession.UserID {
		t.Errorf("UserID mismatch. Expected: %d, Got: %d", originalSession.UserID, loadedSession.UserID)
	}
	
	if loadedSession.Username != originalSession.Username {
		t.Errorf("Username mismatch. Expected: %s, Got: %s", originalSession.Username, loadedSession.Username)
	}
	
	if loadedSession.Email != originalSession.Email {
		t.Errorf("Email mismatch. Expected: %s, Got: %s", originalSession.Email, loadedSession.Email)
	}
	
	// Check custom data
	if role, exists := loadedSession.Get("role"); !exists || role != "admin" {
		t.Errorf("Expected role=admin, got %s (exists: %v)", role, exists)
	}
	
	if theme, exists := loadedSession.Get("theme"); !exists || theme != "dark" {
		t.Errorf("Expected theme=dark, got %s (exists: %v)", theme, exists)
	}
}

func TestLoadSessionWithNoCookie(t *testing.T) {
	config := DefaultConfig()
	req := httptest.NewRequest("GET", "/", nil)
	
	session, err := LoadSession(req, config)
	if err != nil {
		t.Errorf("LoadSession should not error when no cookie present: %v", err)
	}
	
	if session == nil {
		t.Fatal("LoadSession should return a new session when no cookie present")
	}
	
	if session.IsAuthenticated() {
		t.Error("New session should not be authenticated")
	}
}

func TestLoadSessionWithInvalidCookie(t *testing.T) {
	config := DefaultConfig()
	req := httptest.NewRequest("GET", "/", nil)
	
	// Add invalid cookie
	invalidCookie := &http.Cookie{
		Name:  config.CookieName,
		Value: "invalid-base64-data!@#",
	}
	req.AddCookie(invalidCookie)
	
	session, err := LoadSession(req, config)
	if err != nil {
		t.Errorf("LoadSession should not error with invalid cookie: %v", err)
	}
	
	if session.IsAuthenticated() {
		t.Error("Session from invalid cookie should not be authenticated")
	}
}

func TestDestroySession(t *testing.T) {
	config := DefaultConfig()
	w := httptest.NewRecorder()
	
	DestroySession(w, config)
	
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies were set for session destruction")
	}
	
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == config.CookieName {
			sessionCookie = cookie
			break
		}
	}
	
	if sessionCookie == nil {
		t.Fatal("Session cookie not found in destroy response")
	}
	
	if sessionCookie.MaxAge != -1 {
		t.Errorf("Expected MaxAge -1 for destroyed cookie, got %d", sessionCookie.MaxAge)
	}
	
	if sessionCookie.Value != "" {
		t.Error("Destroyed cookie should have empty value")
	}
}

func TestContextOperations(t *testing.T) {
	session := &Session{
		UserID:   456,
		Username: "contextuser",
	}
	
	ctx := context.Background()
	
	// Test WithSession
	ctxWithSession := WithSession(ctx, session)
	
	// Test FromContext
	retrievedSession := FromContext(ctxWithSession)
	if retrievedSession == nil {
		t.Fatal("FromContext returned nil")
	}
	
	if retrievedSession.UserID != session.UserID {
		t.Errorf("UserID mismatch from context. Expected: %d, Got: %d", session.UserID, retrievedSession.UserID)
	}
	
	if retrievedSession.Username != session.Username {
		t.Errorf("Username mismatch from context. Expected: %s, Got: %s", session.Username, retrievedSession.Username)
	}
	
	// Test FromContext with no session
	emptySession := FromContext(context.Background())
	if emptySession.IsAuthenticated() {
		t.Error("FromContext should return new session when none in context")
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test GetSession helper
	session := &Session{UserID: 789}
	ctx := WithSession(context.Background(), session)
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	
	retrievedSession := GetSession(req)
	if retrievedSession.UserID != 789 {
		t.Errorf("GetSession helper failed. Expected UserID 789, got %d", retrievedSession.UserID)
	}
	
	// Test RequireAuth middleware
	authMiddleware := RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	
	// Test with unauthenticated session
	w := httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil).WithContext(WithSession(context.Background(), NewSession()))
	authMiddleware(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("RequireAuth should return 401 for unauthenticated user, got %d", w.Code)
	}
	
	// Test with authenticated session
	w = httptest.NewRecorder()
	authSession := &Session{UserID: 123}
	req = httptest.NewRequest("GET", "/", nil).WithContext(WithSession(context.Background(), authSession))
	authMiddleware(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("RequireAuth should return 200 for authenticated user, got %d", w.Code)
	}
	
	if w.Body.String() != "success" {
		t.Errorf("RequireAuth should call next handler, got body: %s", w.Body.String())
	}
}