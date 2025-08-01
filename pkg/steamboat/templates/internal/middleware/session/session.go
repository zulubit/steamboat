package session

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type contextKey string

const SessionContextKey contextKey = "session"

type Session struct {
	UserID   int               `json:"user_id,omitempty"`
	Username string            `json:"username,omitempty"`
	Email    string            `json:"email,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

type Config struct {
	CookieName string
	SecretKey  string
	MaxAge     int
	HttpOnly   bool
	Secure     bool
	SameSite   http.SameSite
	Domain     string
	Path       string
}

func DefaultConfig() *Config {
	secretKey := os.Getenv("SESSION_KEY")
	if secretKey == "" {
		secretKey = generateRandomKey()
	}

	return &Config{
		CookieName: "steamboat_session",
		SecretKey:  secretKey,
		MaxAge:     86400 * 7, // 7 days
		HttpOnly:   true,
		Secure:     false, // Set to true in production with HTTPS
		SameSite:   http.SameSiteLaxMode,
		Domain:     "",
		Path:       "/",
	}
}

func generateRandomKey() string {
	key := make([]byte, 32)
	rand.Read(key)
	return base64.StdEncoding.EncodeToString(key)
}

func NewSession() *Session {
	return &Session{
		Data: make(map[string]string),
	}
}

func (s *Session) Set(key, value string) {
	if s.Data == nil {
		s.Data = make(map[string]string)
	}
	s.Data[key] = value
}

func (s *Session) Get(key string) (string, bool) {
	if s.Data == nil {
		return "", false
	}
	val, exists := s.Data[key]
	return val, exists
}

func (s *Session) Delete(key string) {
	if s.Data != nil {
		delete(s.Data, key)
	}
}

func (s *Session) Clear() {
	s.UserID = 0
	s.Username = ""
	s.Email = ""
	s.Data = make(map[string]string)
}

func (s *Session) IsAuthenticated() bool {
	return s.UserID > 0
}

func SaveSession(w http.ResponseWriter, session *Session, config *Config) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	encrypted, err := encrypt(data, config.SecretKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt session: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(encrypted)

	cookie := &http.Cookie{
		Name:     config.CookieName,
		Value:    encoded,
		MaxAge:   config.MaxAge,
		HttpOnly: config.HttpOnly,
		Secure:   config.Secure,
		SameSite: config.SameSite,
		Domain:   config.Domain,
		Path:     config.Path,
	}

	http.SetCookie(w, cookie)
	return nil
}

func LoadSession(r *http.Request, config *Config) (*Session, error) {
	cookie, err := r.Cookie(config.CookieName)
	if err != nil {
		return NewSession(), nil
	}

	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return NewSession(), nil
	}

	decrypted, err := decrypt(data, config.SecretKey)
	if err != nil {
		return NewSession(), nil
	}

	var session Session
	if err := json.Unmarshal(decrypted, &session); err != nil {
		return NewSession(), nil
	}

	if session.Data == nil {
		session.Data = make(map[string]string)
	}

	return &session, nil
}

func DestroySession(w http.ResponseWriter, config *Config) {
	cookie := &http.Cookie{
		Name:     config.CookieName,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: config.HttpOnly,
		Secure:   config.Secure,
		SameSite: config.SameSite,
		Domain:   config.Domain,
		Path:     config.Path,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
}

func FromContext(ctx context.Context) *Session {
	if session, ok := ctx.Value(SessionContextKey).(*Session); ok {
		return session
	}
	return NewSession()
}

func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, SessionContextKey, session)
}

func deriveKey(secretKey string) []byte {
	hash := sha256.Sum256([]byte(secretKey))
	return hash[:]
}

func encrypt(data []byte, secretKey string) ([]byte, error) {
	key := deriveKey(secretKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, secretKey string) ([]byte, error) {
	key := deriveKey(secretKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, fmt.Errorf("malformed ciphertext")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

//Helpers

func GetSession(r *http.Request) *Session {
	return FromContext(r.Context())
}

func SaveSessionHelper(w http.ResponseWriter, sess *Session, config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}
	return SaveSession(w, sess, config)
}

func DestroySessionHelper(w http.ResponseWriter, config *Config) {
	if config == nil {
		config = DefaultConfig()
	}
	DestroySession(w, config)
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := GetSession(r)
		if !sess.IsAuthenticated() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
