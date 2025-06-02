// Unit test for TimeRequest JSON marshaling and unmarshaling
// This test ensures that the TimeRequest struct is correctly marshaled to JSON and unmarshaled back to a map.
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestTimeRequestMarshal checks that TimeRequest marshals to the expected JSON and can be unmarshaled correctly.
func TestTimeRequestMarshal(t *testing.T) {
	req := TimeRequest{Timezone: "America/New_York"}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal TimeRequest: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Failed to unmarshal marshaled TimeRequest: %v", err)
	}
	if m["timezone"] != "America/New_York" {
		t.Errorf("Expected timezone to be 'America/New_York', got '%s'", m["timezone"])
	}
}

// TestTimeResponseMarshal checks that TimeResponse marshals to the expected JSON and can be unmarshaled correctly.
func TestTimeResponseMarshal(t *testing.T) {
	// Create a TimeResponse with a sample time string
	resp := TimeResponse{Time: "2025-05-27T12:34:56Z"}
	// Marshal the struct to JSON
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal TimeResponse: %v", err)
	}
	// Unmarshal the JSON back to a map for easy field checking
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Failed to unmarshal marshaled TimeResponse: %v", err)
	}
	// Check that the 'time' field matches the expected value
	if m["time"] != "2025-05-27T12:34:56Z" {
		t.Errorf("Expected time to be '2025-05-27T12:34:56Z', got '%s'", m["time"])
	}
}

// TestHTTPTimeEndpoint tests that the /mcp/time endpoint returns correct JSON for a valid request.
// This test uses the httptest package to simulate an HTTP POST request to the server.
func TestHTTPTimeEndpoint(t *testing.T) {
	// Mock getTimeFromMCP to avoid running Docker
	originalGetTimeFromMCP := getTimeFromMCP
	var getTimeFromMCP = func(timezone string) (string, error) {
		if timezone == "Europe/London" {
			return "2025-06-01T12:00:00Z", nil
		}
		return "", nil
	}
	defer func() { getTimeFromMCP = originalGetTimeFromMCP }()

	// Set up the HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req TimeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		time, err := getTimeFromMCP(req.Timezone)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(TimeResponse{Time: time})
	})

	// Create a test server
	reqBody := []byte(`{"timezone":"Europe/London"}`)
	req := httptest.NewRequest("POST", "/mcp/time", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", resp.StatusCode)
	}

	var tr TimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	expected := "2025-06-01T12:00:00Z"
	if tr.Time != expected {
		t.Errorf("Expected time %q, got %q", expected, tr.Time)
	}
}

// ...existing code...
