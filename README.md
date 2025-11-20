Steps to Run
------------

* Export Env variables
  * export DATABASE_DSN="file:employee.db?cache=shared&_fk=1"
  * export SERVER_ADDR=":8080"
* go run ./cmd/server/main.go

---

# Employee Management REST API

A production-ready RESTful API service for managing employee records, built with Go and following clean architecture principles.

## Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [Testing](#testing)
- [Development](#development)
- [Deployment](#deployment)

## Overview

This application provides a complete employee management system with RESTful APIs for creating, reading, updating, and deleting employee records. It's built using modern Go practices with a layered architecture that separates concerns and promotes maintainability.

**Key Technologies:**
- **Language:** Go 1.25.1
- **Router:** Chi v5 (lightweight, idiomatic HTTP router)
- **Database:** SQLite (with support for PostgreSQL via DSN configuration)
- **ORM/Query Builder:** sqlx (provides extensions on top of database/sql)
- **Testing:** testify, go-sqlmock

## Architecture

The application follows a **3-tier layered architecture** with clear separation of concerns:

```
┌─────────────────────────────────────┐
│         HTTP Handler Layer          │  ← Handles HTTP requests/responses
│     (internal/handler)              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│        Business Logic Layer         │  ← Contains business rules
│     (internal/service)              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│       Data Access Layer (DAO)       │  ← Database operations
│     (internal/dao)                  │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│          Database (SQLite)          │
└─────────────────────────────────────┘
```

### Layer Responsibilities

1. **Handler Layer** ([internal/handler/employee_handler.go](internal/handler/employee_handler.go))
   - Receives HTTP requests
   - Validates and deserializes JSON
   - Calls service layer
   - Serializes responses
   - Returns appropriate HTTP status codes

2. **Service Layer** ([internal/service/employee_service.go](internal/service/employee_service.go))
   - Implements business logic
   - Validates business rules
   - Orchestrates DAO operations
   - Returns domain-specific errors

3. **DAO Layer** ([internal/dao/employee_dao.go](internal/dao/employee_dao.go))
   - Direct database interactions
   - CRUD operations
   - Query execution
   - Transaction management

## Features

- **CRUD Operations:** Complete Create, Read, Update, Delete functionality for employees
- **RESTful API Design:** Follows REST conventions with proper HTTP methods and status codes
- **Graceful Shutdown:** Properly handles shutdown signals to complete in-flight requests
- **Connection Pooling:** Configurable database connection pool for optimal performance
- **Middleware Stack:**
  - Request ID tracking for debugging
  - Request/response logging
  - Panic recovery
  - Request timeout protection (30s)
- **Health Check Endpoint:** Monitor application status
- **Clean Architecture:** Easily testable with mock interfaces
- **Environment-based Configuration:** Flexible configuration via environment variables

## Prerequisites

- **Go 1.25.1 or higher**
- **SQLite3** (installed by default on most systems)
- **Git** (for cloning the repository)

## Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd employee-app-go
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build the application:**
   ```bash
   go build -o bin/server ./cmd/server
   ```

## Configuration

The application is configured through environment variables with sensible defaults:

| Environment Variable | Description | Default Value |
|---------------------|-------------|---------------|
| `SERVER_ADDR` | HTTP server address and port | `:8080` |
| `DATABASE_DSN` | Database connection string | `file:employees.db?_busy_timeout=5000&_foreign_keys=1` |
| `DB_MAX_OPEN_CONNS` | Maximum open database connections | `25` |
| `DB_MAX_IDLE_CONNS` | Maximum idle database connections | `25` |
| `DB_CONN_MAX_LIFETIME_SECONDS` | Connection max lifetime in seconds | `300` |

### Configuration Examples

**Development (SQLite):**
```bash
export DATABASE_DSN="file:employees.db?_busy_timeout=5000&_foreign_keys=1"
export SERVER_ADDR=":8080"
```

**Production (PostgreSQL):**
```bash
export DATABASE_DSN="postgres://user:password@localhost:5432/employees?sslmode=require"
export SERVER_ADDR=":8080"
export DB_MAX_OPEN_CONNS="100"
export DB_MAX_IDLE_CONNS="50"
```

## API Endpoints

Base URL: `http://localhost:8080`

### Employee Endpoints

| Method | Endpoint | Description | Request Body | Response |
|--------|----------|-------------|--------------|----------|
| `POST` | `/api/v1/employees/` | Create new employee | Employee JSON | 201 Created + Employee object |
| `GET` | `/api/v1/employees/` | List all employees | - | 200 OK + Employee array |
| `GET` | `/api/v1/employees/{id}/` | Get employee by ID | - | 200 OK + Employee object |
| `PUT` | `/api/v1/employees/{id}/` | Update employee | Employee JSON | 200 OK + Updated employee |
| `DELETE` | `/api/v1/employees/{id}/` | Delete employee | - | 204 No Content |

### Health Check

| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| `GET` | `/health` | Health check | 200 OK + "ok" |

### Request/Response Examples

**Create Employee:**
```bash
curl -X POST http://localhost:8080/api/v1/employees/ \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "position": "Software Engineer"
  }'
```

Response (201 Created):
```json
{
  "id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "position": "Software Engineer",
  "created_at": "2024-11-20T10:30:00Z",
  "updated_at": "2024-11-20T10:30:00Z"
}
```

**Get All Employees:**
```bash
curl http://localhost:8080/api/v1/employees/
```

**Get Employee by ID:**
```bash
curl http://localhost:8080/api/v1/employees/1/
```

**Update Employee:**
```bash
curl -X PUT http://localhost:8080/api/v1/employees/1/ \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "position": "Senior Software Engineer"
  }'
```

**Delete Employee:**
```bash
curl -X DELETE http://localhost:8080/api/v1/employees/1/
```

## Project Structure

```
employee-app-go/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading and defaults
│   ├── db/
│   │   └── pool.go              # Database connection pool setup
│   ├── model/
│   │   └── model.go             # Employee data model
│   ├── dao/
│   │   ├── employee_dao.go      # Data Access Object interface & impl
│   │   └── employee_dao_test.go # DAO unit tests with mocks
│   ├── service/
│   │   └── employee_service.go  # Business logic layer
│   ├── handler/
│   │   └── employee_handler.go  # HTTP handlers
│   └── router/
│       └── router.go            # Route definitions and middleware
├── employees.db                 # SQLite database file (auto-created)
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
└── README.md                    # This file
```

### Package Descriptions

- **[cmd/server/main.go](cmd/server/main.go)**: Application bootstrap, dependency injection, server lifecycle
- **[internal/config](internal/config/config.go)**: Centralized configuration with environment variable support
- **[internal/db](internal/db/pool.go)**: Database connection management, pooling, and schema initialization
- **[internal/model](internal/model/model.go)**: Employee struct with JSON and database tags
- **[internal/dao](internal/dao/employee_dao.go)**: Database operations (Create, Read, Update, Delete)
- **[internal/service](internal/service/employee_service.go)**: Business logic, validation, error handling
- **[internal/handler](internal/handler/employee_handler.go)**: HTTP request/response handling
- **[internal/router](internal/router/router.go)**: Route mapping and middleware configuration

## Database Schema

The application automatically creates the required schema on startup.

**Employees Table:**
```sql
CREATE TABLE IF NOT EXISTS employees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    position TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Fields:**
- `id`: Auto-incrementing primary key
- `first_name`: Employee's first name (required)
- `last_name`: Employee's last name (required)
- `email`: Unique email address (required)
- `position`: Job position/title (optional)
- `created_at`: Record creation timestamp (auto-generated)
- `updated_at`: Last update timestamp (auto-updated)

## Testing

The project includes comprehensive unit tests with mocking.

**Run all tests:**
```bash
go test ./...
```

**Run tests with coverage:**
```bash
go test -cover ./...
```

**Run tests with verbose output:**
```bash
go test -v ./...
```

**Run tests for a specific package:**
```bash
go test ./internal/dao/...
```

**Test Coverage Report:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Files
- [internal/dao/employee_dao_test.go](internal/dao/employee_dao_test.go): DAO layer tests using sqlmock

## Development

### Running in Development Mode

1. **Set environment variables:**
   ```bash
   export DATABASE_DSN="file:employees.db?_busy_timeout=5000&_foreign_keys=1"
   export SERVER_ADDR=":8080"
   ```

2. **Run the server:**
   ```bash
   go run ./cmd/server/main.go
   ```

3. **Access the API:**
   - API Base: `http://localhost:8080/api/v1`
   - Health Check: `http://localhost:8080/health`

### Adding New Features

To add a new entity (e.g., Department):

1. **Define model** in `internal/model/department.go`
2. **Create DAO interface and implementation** in `internal/dao/department_dao.go`
3. **Write DAO tests** in `internal/dao/department_dao_test.go`
4. **Implement service layer** in `internal/service/department_service.go`
5. **Create HTTP handlers** in `internal/handler/department_handler.go`
6. **Add routes** in `internal/router/router.go`
7. **Update database schema** in `internal/db/pool.go`

### Code Style

- Follow standard Go conventions (gofmt, golint)
- Use meaningful variable names
- Add comments for exported functions and types
- Keep functions small and focused
- Write tests for all business logic

## Deployment

### Building for Production

```bash
# Build for current platform
go build -o bin/server ./cmd/server

# Build for Linux (common for cloud deployment)
GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/server
```

### Docker Deployment

Create a `Dockerfile`:
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

Build and run:
```bash
docker build -t employee-app .
docker run -p 8080:8080 \
  -e DATABASE_DSN="file:employees.db?_busy_timeout=5000&_foreign_keys=1" \
  employee-app
```

### Environment Setup

**Production Checklist:**
- [ ] Use PostgreSQL instead of SQLite for multi-user scenarios
- [ ] Configure proper connection pool sizes
- [ ] Enable HTTPS/TLS
- [ ] Set up logging to external service
- [ ] Configure health check monitoring
- [ ] Set appropriate timeouts
- [ ] Use secrets management for sensitive config
- [ ] Enable CORS if needed for web clients
- [ ] Set up database backups

## Dependencies

This project uses the following key dependencies:

- **[chi/v5](https://github.com/go-chi/chi)**: Lightweight, composable HTTP router
- **[sqlx](https://github.com/jmoiron/sqlx)**: Extensions to Go's database/sql package
- **[go-sqlite3](https://github.com/mattn/go-sqlite3)**: SQLite driver for Go
- **[testify](https://github.com/stretchr/testify)**: Testing toolkit with assertions
- **[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)**: Mock SQL driver for testing

## Troubleshooting

**Port already in use:**
```bash
# Change the server port
export SERVER_ADDR=":8081"
```

**Database locked:**
```bash
# Increase busy timeout in DSN
export DATABASE_DSN="file:employees.db?_busy_timeout=10000&_foreign_keys=1"
```

**Module issues:**
```bash
# Clean and rebuild
go clean -modcache
go mod tidy
go mod download
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Contact

For questions or support, please open an issue in the repository.
