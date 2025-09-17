# CGC-2026 Server

A Go-based HTTP server with basic health endpoint.

## Prerequisites

- Go 1.25 or higher: https://go.dev/dl/

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

## Using Make

This project includes a Makefile with common commands:

```bash
make build    # Build the application
make run      # Run the application
make dev      # Run the application with live reload
make test     # Run tests
make clean    # Clean build artifacts
```

## Live Reload for Development

This project supports live reload using [Air](https://github.com/air-verse/air) via Docker. This allows you to automatically rebuild and restart the server when code changes are detected.

### Using Live Reload

To use live reload, simply run:

```bash
# Using the make command (recommended)
make dev

# Or directly with docker-compose
docker-compose up dev
```

The configuration for Air is in the `.air.toml` file in the project root. No need to install Air locally as it runs in a Docker container.

## Project Structure

- `/cmd/server`: Entry point for the application
- `/internal`: Internal packages not meant for external use
  - `/app`: Core application logic
    - `/handlers`: Feature-specific HTTP handlers
    - `/router`: Route registration and organization
  - `/log`: Logging utilities
  - `/http`: HTTP utilities and handlers
    - `/middleware`: HTTP middleware components
    - `/responses`: Standardized HTTP response helpers
- `/api`: API definitions and documentation