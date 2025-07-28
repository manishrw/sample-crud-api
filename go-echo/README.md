# Go Echo CRUD API

A RESTful CRUD API built with Go Echo framework following industry best practices.

## Features

- **RESTful API**: Complete CRUD operations for users
- **Clean Architecture**: Separation of concerns with handlers, services, and repositories
- **Go Best Practices**: Internal package structure for encapsulation and module boundaries
- **Database Integration**: PostgreSQL with proper connection pooling
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Validation**: Input validation and business logic validation
- **Graceful Shutdown**: Proper server shutdown handling
- **Configuration Management**: Viper-based configuration with support for YAML files and environment variables
- **CORS Support**: Cross-origin resource sharing enabled
- **Health Check**: Built-in health check endpoint
- **Comprehensive Testing**: Unit tests for all layers with mocking

## Requirements

- Go 1.24 or higher
- PostgreSQL 12 or higher

## Package Structure

The application follows Go best practices with an `internal/` package structure:

### Internal Package (`internal/`)
The `internal/` directory contains all application-specific code that should not be imported by other modules. This follows Go's convention for package encapsulation:

- **`internal/models/`**: Data structures and DTOs
- **`internal/repository/`**: Data access layer with database operations
- **`internal/services/`**: Business logic and validation
- **`internal/handlers/`**: HTTP request/response handling

### Public Packages
- **`config/`**: Configuration management (can be imported by other modules)
- **`database/`**: Database connection utilities (can be imported by other modules)

This structure ensures that:
- Internal business logic is protected from external imports
- The module maintains clear boundaries
- Only necessary packages are exposed for external use
- The codebase follows Go community conventions

## Project Structure

```
go-echo/
├── config/                 # Configuration management
│   ├── config.go          # Viper-based configuration
│   ├── config.yaml        # Sample configuration file
│   └── config_test.go     # Configuration tests
├── database/              # Database connection and migrations
│   ├── database.go
│   └── migrations/
│       └── 001_create_users_table.sql
├── internal/              # Internal application code (not importable)
│   ├── handlers/          # HTTP request handlers
│   │   ├── user_handler.go
│   │   └── user_handler_test.go
│   ├── models/            # Data models and DTOs
│   │   └── user.go
│   ├── repository/        # Data access layer
│   │   ├── user_repository.go
│   │   └── user_repository_test.go
│   └── services/          # Business logic layer
│       ├── user_service.go
│       └── user_service_test.go
├── scripts/               # Helper scripts
│   └── docker-build.sh    # Docker build and run script
├── .dockerignore          # Docker ignore file
├── config.yaml            # Root configuration file
├── config.yaml.example    # Example configuration file
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile             # Docker build file
├── go.mod                 # Go module file
├── main.go                # Application entry point
├── Makefile               # Development commands
└── README.md              # This file
```

## Configuration

The application uses Viper for configuration management, supporting multiple configuration sources:

### Configuration Sources (in order of precedence):
1. **Environment Variables** (highest priority)
2. **Configuration Files** (YAML format)
3. **Default Values** (lowest priority)

### Configuration File

Create a `config.yaml` file in the project root or `config/` directory:

```yaml
server:
  host: "localhost"
  port: "8080"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  name: "users_db"
  sslmode: "disable"

jwt:
  secret: "your-secret-key-here"
  expiry_hours: 24
```

### Environment Variables

You can override any configuration using environment variables:

```bash
export SERVER_HOST="0.0.0.0"
export SERVER_PORT="9090"
export DATABASE_HOST="db.example.com"
export DATABASE_USER="myuser"
export DATABASE_PASSWORD="mypassword"
export DATABASE_NAME="mydb"
export JWT_SECRET="my-secret-key"
export JWT_EXPIRY_HOURS="48"
```

### Environment Variable Mapping

| Config Path | Environment Variable | Default Value |
|-------------|---------------------|---------------|
| `server.host` | `SERVER_HOST` | `localhost` |
| `server.port` | `SERVER_PORT` | `8080` |
| `database.host` | `DATABASE_HOST` | `localhost` |
| `database.port` | `DATABASE_PORT` | `5432` |
| `database.user` | `DATABASE_USER` | `postgres` |
| `database.password` | `DATABASE_PASSWORD` | `password` |
| `database.name` | `DATABASE_NAME` | `users_db` |
| `database.sslmode` | `DATABASE_SSLMODE` | `disable` |
| `jwt.secret` | `JWT_SECRET` | `your-secret-key-here` |
| `jwt.expiry_hours` | `JWT_EXPIRY_HOURS` | `24` |

## API Endpoints

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | Get all users |
| GET | `/api/v1/users/:id` | Get user by ID |
| POST | `/api/v1/users` | Create a new user |
| PUT | `/api/v1/users/:id` | Update user by ID |
| DELETE | `/api/v1/users/:id` | Delete user by ID |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |

## Request/Response Examples

### Create User

**Request:**
```bash
POST /api/v1/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

**Response:**
```json
{
  "message": "User created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### Get Users

**Request:**
```bash
GET /api/v1/users
```

**Response:**
```json
{
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

## Installation and Setup

### Quick Start (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd go-echo

# Quick start with Docker Compose
make quick-start
```

### Manual Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd go-echo
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up the database:**
   ```bash
   # Create PostgreSQL database
   createdb users_db
   
   # Run migrations
   make db-migrate
   ```

4. **Configure the application:**
   ```bash
   # Copy and edit the configuration file
   make dev-setup
   # Edit config.yaml with your settings
   ```

5. **Run the application:**
   ```bash
   make run
   ```

### Using Makefile

The project includes a Makefile with common commands:

```bash
# View all available commands
make help

# Build the application
make build

# Run locally
make run

# Run tests
make test

# Run tests with coverage
make test-cover

# Clean build artifacts
make clean

# Docker commands
make docker-build
make docker-compose
make docker-clean
make docker-logs
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./handlers
go test ./services
go test ./repository
go test ./config
```

## Docker Support

The application includes Docker support for easy deployment with Go 1.24:

### Using Docker Compose (Recommended)

```bash
# Build and run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Using the Helper Script

For convenience, you can use the provided helper script:

```bash
# Make script executable (first time only)
chmod +x scripts/docker-build.sh

# Build and run (auto-detects Docker Compose)
./scripts/docker-build.sh run

# Build and run with Docker Compose
./scripts/docker-build.sh compose

# Build and run with Docker directly
./scripts/docker-build.sh docker

# Clean up Docker resources
./scripts/docker-build.sh clean
```

### Using Docker directly

```bash
# Build the Docker image
docker build -t go-echo-api .

# Run the container
docker run -p 8080:8080 \
  -e DATABASE_HOST=your-db-host \
  -e DATABASE_USER=your-db-user \
  -e DATABASE_PASSWORD=your-db-password \
  -e DATABASE_NAME=your-db-name \
  go-echo-api
```

### Docker Configuration

The Docker setup includes:
- **Go 1.24**: Latest Go version for optimal performance
- **Multi-stage build**: Optimized image size
- **Alpine Linux**: Lightweight base image
- **Configuration support**: Environment variables and config files
- **Database migrations**: Automatic migration support

## Development

### Code Structure

The application follows clean architecture principles and Go best practices:

- **internal/**: Internal application code (not importable by other modules)
  - **handlers/**: HTTP request/response handling
  - **services/**: Business logic and validation
  - **repository/**: Data access layer
  - **models/**: Data structures and DTOs
- **config/**: Configuration management with Viper
- **database/**: Database connection and migrations

### Adding New Features

1. Define the model in `internal/models/`
2. Create repository methods in `internal/repository/`
3. Implement business logic in `internal/services/`
4. Add HTTP handlers in `internal/handlers/`
5. Write tests for each layer
6. Update configuration if needed

## License

This project is licensed under the MIT License. 