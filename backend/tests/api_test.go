package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"wait-to-go/auth"
)

type mockApp struct {
	db      *mockDB
	queue   *[]Entry
	history *[]Entry
}

type mockDB struct {
	entries map[int]Entry
	nextID  int
}

type Entry struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Status      string    `json:"status"`
	JoinTime    time.Time `json:"joinTime"`
}

func newMockApp() *mockApp {
	return &mockApp{
		db: &mockDB{
			entries: make(map[int]Entry),
			nextID:  1,
		},
		queue:   &[]Entry{},
		history: &[]Entry{},
	}
}

func (a *mockApp) handleJoin(w http.ResponseWriter, r *http.Request) {
	var entry Entry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if entry.FirstName == "" || entry.LastName == "" || entry.PhoneNumber == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	entry.ID = a.db.nextID
	a.db.nextID++
	entry.Status = "waiting"
	entry.JoinTime = time.Now()

	a.db.entries[entry.ID] = entry
	*a.queue = append(*a.queue, entry)

	token, _ := auth.GenerateToken(entry.ID, entry.PhoneNumber)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    entry.ID,
		"token": token,
	})
}

func (a *mockApp) handleStatus(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetClaimsFromContext(r.Context())
	entry, exists := a.db.entries[claims.ID]
	if !exists {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	position := 0
	if entry.Status == "waiting" {
		for _, e := range *a.queue {
			if e.Status == "waiting" && e.JoinTime.Before(entry.JoinTime) {
				position++
			}
		}
		position++ // Add 1 because we want 1-based position
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"entry":    entry,
		"position": position,
	})
}

func (a *mockApp) handleQueue(w http.ResponseWriter, r *http.Request) {
	var waitingEntries []Entry
	for _, entry := range *a.queue {
		if entry.Status == "waiting" {
			waitingEntries = append(waitingEntries, entry)
		}
	}
	json.NewEncoder(w).Encode(waitingEntries)
}

func TestHandleJoin(t *testing.T) {
	app := newMockApp()

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid join request",
			payload: map[string]interface{}{
				"firstName":   "John",
				"lastName":    "Doe",
				"email":       "john@example.com",
				"phoneNumber": "1234567890",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing required fields",
			payload: map[string]interface{}{
				"firstName": "John",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/join", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			app.handleJoin(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				json.Unmarshal(rr.Body.Bytes(), &response)

				if _, ok := response["token"]; !ok {
					t.Error("Response missing token")
				}
				if _, ok := response["id"]; !ok {
					t.Error("Response missing id")
				}
			}
		})
	}
}

func TestHandleStatus(t *testing.T) {
	app := newMockApp()

	// Add a test entry
	testEntry := Entry{
		ID:          1,
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Status:      "waiting",
		JoinTime:    time.Now(),
	}
	app.db.entries[1] = testEntry
	*app.queue = append(*app.queue, testEntry)

	// Generate a valid token for ID 1
	token, _ := auth.GenerateToken(1, "1234567890")
	// Generate a valid token for non-existent ID
	invalidToken, _ := auth.GenerateToken(999, "9999999999")

	tests := []struct {
		name           string
		id             string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid request",
			id:             "1",
			token:          "Bearer " + token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID",
			id:             "999",
			token:          "Bearer " + invalidToken,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Missing token",
			id:             "1",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/status/"+tt.id, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			rr := httptest.NewRecorder()

			handler := auth.AuthMiddleware(app.handleStatus)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response struct {
					Entry    Entry `json:"entry"`
					Position int   `json:"position"`
				}
				json.Unmarshal(rr.Body.Bytes(), &response)

				if response.Entry.ID != 1 {
					t.Errorf("Expected entry ID 1, got %d", response.Entry.ID)
				}
				if response.Position != 1 {
					t.Errorf("Expected position 1, got %d", response.Position)
				}
			}
		})
	}
}

func TestHandleQueue(t *testing.T) {
	app := newMockApp()

	// Add some test entries
	entries := []Entry{
		{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Status:    "waiting",
			JoinTime:  time.Now(),
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
			Status:    "waiting",
			JoinTime:  time.Now(),
		},
	}

	for _, entry := range entries {
		app.db.entries[entry.ID] = entry
		*app.queue = append(*app.queue, entry)
	}

	// Initialize admin key for testing
	testAdminKey := "test-admin-key"
	auth.AddAdminKey(testAdminKey)

	for _, tt := range []struct {
		name           string
		apiKey         string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Valid request",
			apiKey:         testAdminKey,
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Missing API key",
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
			expectedCount:  0,
		},
		{
			name:           "Invalid API key",
			apiKey:         "invalid-key",
			expectedStatus: http.StatusUnauthorized,
			expectedCount:  0,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/queue", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			rr := httptest.NewRecorder()

			handler := auth.AdminAuthMiddleware(app.handleQueue)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response []Entry
				json.Unmarshal(rr.Body.Bytes(), &response)

				if len(response) != tt.expectedCount {
					t.Errorf("Expected %d entries, got %d", tt.expectedCount, len(response))
				}
			}
		})
	}
}

func TestAdminAuthMiddleware(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Initialize admin key for testing
	testAdminKey := "test-admin-key"
	auth.AddAdminKey(testAdminKey)

	for _, tt := range []struct {
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
	} {
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
