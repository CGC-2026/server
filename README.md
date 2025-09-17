# CGC-2026 Server

A Go-based HTTP server with basic health endpoint.

## Prerequisites

- Go 1.25 or higher

## Running the Server

To run the server, navigate to the project root and execute:

```bash
# From the server directory
go run cmd/server/main.go
```

The server will start on port 8080 by default.

## Testing the Health Endpoint

You can test the health endpoint using curl:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "OK",
  "timestamp": "2025-09-17T12:34:56Z"
}
```

## Package Management

This project uses Go Modules for dependency management. Here's how to work with packages:

### Adding a New Package

To add a new dependency to the project:

```bash
go get github.com/example/package
```

This will automatically update your `go.mod` and `go.sum` files.

### Installing Dependencies

After cloning the repository, install all dependencies:

```bash
go mod download
```

### Updating Dependencies

To update dependencies and clean up unused ones:

```bash
go mod tidy
```

This command ensures your `go.mod` file correctly reflects all dependencies used in the codebase.

### Building the Application

To build an executable:

```bash
go build -o server cmd/server/main.go
```

## API Documentation

The API is documented using OpenAPI specification available at `api/openaapi/openapi.yml`.

## Project Structure

- `/cmd/server`: Entry point for the application
- `/internal`: Internal packages not meant for external use
  - `/app`: Core application logic
  - `/log`: Logging utilities
  - `/http`: HTTP utilities and handlers
    - `/middleware`: HTTP middleware components
    - `/responses`: Standardized HTTP response helpers
- `/api`: API definitions and documentation