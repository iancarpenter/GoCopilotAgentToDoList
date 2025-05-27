# GoCopilotAgentToDoList

## Overview
GoCopilotAgentToDoList is a Go-based HTTP server that provides an API endpoint for retrieving the current time in a specified timezone. The server interacts with a Dockerized Model Context Protocol (MCP) service to fetch the time, making it a practical example of integrating Go, Docker, and external services.

## Features
- **/mcp/time API endpoint**: Accepts a POST request with a timezone and returns the current time for that timezone.
- **Docker integration**: Uses Docker to run the MCP service for time retrieval.
- **JSON-based communication**: Both requests and responses use JSON.
- **Unit tests**: Includes basic unit tests for request and response marshaling.

## Prerequisites
- [Go](https://golang.org/dl/) 1.18 or later
- [Docker](https://www.docker.com/products/docker-desktop) installed and running
- MCP Docker image (`mcp:latest`) available locally
- Windows OS (but code is cross-platform)

## Project Structure
```
go.mod                  # Go module definition
package main.go         # Main application source code
package_main_test.go    # Unit tests for the project
todos.json              # (Reserved for future to-do list features)
```

## Usage

### 1. Build and Run the Server
Open a terminal in the project directory and run:

```
go run package\ main.go
```

The server will start and listen on `http://localhost:8080`.

### 2. API Endpoint
#### POST `/mcp/time`
- **Request Body:**
  ```json
  {
    "timezone": "America/New_York"
  }
  ```
- **Response Body:**
  ```json
  {
    "time": "2025-05-27T12:34:56Z"
  }
  ```
- **Example with curl:**
  ```sh
  curl -X POST http://localhost:8080/mcp/time -H "Content-Type: application/json" -d '{"timezone":"America/New_York"}'
  ```

### 3. How It Works
- The server receives a POST request at `/mcp/time` with a JSON body specifying the timezone.
- It marshals the request and pipes it to a Docker container running the MCP service (`mcp:latest`), specifically the `mcp/time` command.
- The output from the container is parsed and returned as a JSON response.

## Running Tests
To run the unit tests, execute:

```
go test -v
```

## Extending the Project
- **To-Do List Feature:** The presence of `todos.json` suggests future plans for a to-do list API.
- **Additional Endpoints:** You can add more endpoints by defining new handlers in `package main.go`.
- **Error Handling:** Improve error messages and logging for production use.

## Troubleshooting
- Ensure Docker is running and the `mcp:latest` image is available locally.
- If you encounter permission issues with Docker, run your terminal as administrator.
- The server listens on port 8080 by default. Change this in `package main.go` if needed.

## License
This project is provided as-is for educational and demonstration purposes.
