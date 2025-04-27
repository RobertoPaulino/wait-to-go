package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-256-bit-secret"))

// AdminKey represents a hashed API key with metadata
type AdminKey struct {
	HashedKey []byte
	CreatedAt time.Time
	LastUsed  time.Time
}

type AdminKeyStore struct {
	keys map[string]*AdminKey
	mu   sync.RWMutex
}

var adminKeyStore = &AdminKeyStore{
	keys: make(map[string]*AdminKey),
}

func init() {
	// Initialize with default admin key if provided in env
	if key := os.Getenv("ADMIN_API_KEY"); key != "" {
		AddAdminKey(key)
	}
}

// AddAdminKey adds a new admin key to the store
func AddAdminKey(key string) error {
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminKeyStore.mu.Lock()
	defer adminKeyStore.mu.Unlock()

	// Use base64 of hashed key as identifier
	keyID := base64.StdEncoding.EncodeToString(hashedKey)
	adminKeyStore.keys[keyID] = &AdminKey{
		HashedKey: hashedKey,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	return nil
}

// RemoveAdminKey removes an admin key from the store
func RemoveAdminKey(keyID string) {
	adminKeyStore.mu.Lock()
	defer adminKeyStore.mu.Unlock()
	delete(adminKeyStore.keys, keyID)
}

// ValidateAdminKey checks if the provided key is valid
func ValidateAdminKey(key string) bool {
	adminKeyStore.mu.RLock()
	defer adminKeyStore.mu.RUnlock()

	// Try each stored key
	for keyID, adminKey := range adminKeyStore.keys {
		if err := bcrypt.CompareHashAndPassword(adminKey.HashedKey, []byte(key)); err == nil {
			// Update last used time
			adminKeyStore.keys[keyID].LastUsed = time.Now()
			return true
		}
	}

	return false
}

// RateLimiter implements a simple token bucket algorithm
type RateLimiter struct {
	tokens     map[string][]time.Time
	windowSize time.Duration
	maxTokens  int
	mu         sync.Mutex
}

func NewRateLimiter(windowSize time.Duration, maxTokens int) *RateLimiter {
	return &RateLimiter{
		tokens:     make(map[string][]time.Time),
		windowSize: windowSize,
		maxTokens:  maxTokens,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if _, exists := rl.tokens[key]; !exists {
		rl.tokens[key] = []time.Time{now}
		return true
	}

	// Remove tokens outside the window
	windowStart := now.Add(-rl.windowSize)
	var validTokens []time.Time
	for _, t := range rl.tokens[key] {
		if t.After(windowStart) {
			validTokens = append(validTokens, t)
		}
	}

	if len(validTokens) >= rl.maxTokens {
		return false
	}

	rl.tokens[key] = append(validTokens, now)
	return true
}

// Initialize rate limiters
var (
	authLimiter  = NewRateLimiter(time.Minute, 30)  // 30 requests per minute for auth
	adminLimiter = NewRateLimiter(time.Minute, 100) // 100 requests per minute for admin
)

type Claims struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

func GenerateToken(id int, phoneNumber string) (string, error) {
	claims := Claims{
		ID:          id,
		PhoneNumber: phoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Rate limiting based on IP
		clientIP := r.RemoteAddr
		if !authLimiter.Allow(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		r = r.WithContext(AddClaimsToContext(r.Context(), claims))
		next(w, r)
	}
}

func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Rate limiting based on IP
		clientIP := r.RemoteAddr
		if !adminLimiter.Allow(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		if !ValidateAdminKey(apiKey) {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
