package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
)

// TimeRequest and TimeResponse are the expected request/response types for the mcp/time api
type TimeRequest struct {
	Timezone string `json:"timezone"`
}
type TimeResponse struct {
	Time string `json:"time"`
}

func getTimeFromMCP(timezone string) (string, error) {
	// Prepare the request
	req := TimeRequest{Timezone: timezone}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// run the docker command
	cmd := exec.Command("docker", "run", "--rm", "-i", "mcp:latest", "mcp/time")
	cmd.Stdin = bytes.NewReader(reqBytes)
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", err
	}
	// Parse the response
	var res TimeResponse
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		return "", err
	}
	return res.Time, nil
}

// main is the entry point of the application. It sets up HTTP handlers and starts the server.
func main() {
	// Set up the HTTP server and handlers
	http.HandleFunc("/mcp/time", func(w http.ResponseWriter, r *http.Request) {
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

	http.ListenAndServe(":8080", nil)
}
