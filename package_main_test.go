// Unit test for TimeRequest JSON marshaling and unmarshaling
// This test ensures that the TimeRequest struct is correctly marshaled to JSON and unmarshaled back to a map.
package main

import (
	"encoding/json"
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
