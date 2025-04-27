package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"wait-to-go/auth"
)

func TestGenerateAndValidateToken(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		phoneNumber string
		wantErr     bool
	}{
		{
			name:        "Valid token",
			id:          1,
			phoneNumber: "1234567890",
			wantErr:     false,
		},
		{
			name:        "Empty phone number",
			id:          1,
			phoneNumber: "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.GenerateToken(tt.id, tt.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			claims, err := auth.ValidateToken(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if claims.ID != tt.id {
					t.Errorf("ValidateToken() got ID = %v, want %v", claims.ID, tt.id)
				}
				if claims.PhoneNumber != tt.phoneNumber {
					t.Errorf("ValidateToken() got PhoneNumber = %v, want %v", claims.PhoneNumber, tt.phoneNumber)
				}
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	// Generate a valid token for testing
	token, _ := auth.GenerateToken(1, "1234567890")

	// Create a simple handler for testing
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			token:          "Bearer " + token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token format",
			token:          "Invalid " + token,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			rr := httptest.NewRecorder()

			handler := auth.AuthMiddleware(testHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestAdminKeyValidation(t *testing.T) {
	// Initialize admin key for testing
	testAdminKey := "test-admin-key-2"
	if err := auth.AddAdminKey(testAdminKey); err != nil {
		t.Fatalf("Failed to add admin key: %v", err)
	}

	// Create a simple handler for testing
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		apiKey         string
		expectedStatus int
	}{
		{
			name:           "Valid API key",
			apiKey:         testAdminKey,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing API key",
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid API key",
			apiKey:         "invalid-key",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			rr := httptest.NewRecorder()

			handler := auth.AdminAuthMiddleware(testHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
