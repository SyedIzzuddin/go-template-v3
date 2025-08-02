# Go REST API Template

A production-ready, scalable Go REST API template built with modern best practices. This template provides a solid foundation for building RESTful APIs with user management, file upload capabilities, and comprehensive tooling for rapid development.

## 🚀 Features

- **Clean Architecture**: Layered architecture with clear separation of concerns (handler → service → repository → entity)
- **Echo Framework**: High-performance HTTP router and middleware
- **PostgreSQL Integration**: Raw SQL with pgx driver and SQLC for type-safe queries
- **Database Migrations**: Goose for schema versioning with automatic migration on startup
- **File Management**: Secure file upload with validation and storage
- **Structured Logging**: Zap logger with request tracing
- **Security Features**: Rate limiting, CORS, input validation
- **Testing**: Integration tests with Testcontainers
- **Live Reload**: Air for development with hot reloading
- **Docker Support**: Complete containerization with Docker Compose

## 🏗️ Architecture

```
├── cmd/api/                    # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── database/               # Database connection and health checks
│   ├── dto/                    # Data Transfer Objects
│   ├── entity/                 # Domain entities (User, File)
│   ├── handler/                # HTTP handlers (controllers)
│   ├── logger/                 # Structured logging configuration
│   ├── middleware/             # Custom middleware (rate limiting, CORS, logging)
│   ├── repository/             # Data access layer with raw SQL
│   ├── router/                 # Route definitions
│   ├── server/                 # Server initialization and configuration
│   ├── service/                # Business logic layer
│   └── utils/                  # Utility functions
├── pkg/
│   ├── response/               # Standardized API responses
│   ├── storage/                # File storage utilities
│   └── validator/              # Request validation
├── db/
│   ├── migrations/             # Goose migration files
│   ├── queries/                # SQL query files for SQLC
│   └── sqlc/                   # Generated SQLC code (auto-generated)
├── uploads/                    # File upload directory
└── docs/                       # API documentation
```

## 🚦 Quick Start

### Prerequisites

- Go 1.23.4 or higher
- Docker and Docker Compose
- Make (for using Makefile commands)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd go-template-v3
cp .env.example .env
```

### 2. Choose Your Development Workflow

We provide multiple ways to run the application. Choose the one that fits your workflow:

#### Option A: Full Docker Setup (Production-like)

**Best for:** Testing complete setup, production simulation

```bash
# Starts both database and application in Docker
make docker-run
```

#### Option B: Hybrid Development (Recommended for Development)

**Best for:** Active development with fast rebuilds and debugging

```bash
# Terminal 1: Start only the database
docker compose up psql_bp -d

# Terminal 2: Run app locally with hot reload
make watch
```

#### Option C: Local Development

**Best for:** Development without hot reload

```bash
# Terminal 1: Start only the database
docker compose up psql_bp -d

# Terminal 2: Run app locally
make run
```

### 3. Setup Database (for Options B & C)

```bash
# Install required tools (sqlc, goose)
make install-tools

# Run database migrations (only needed if DB_AUTO_MIGRATE=false)
DATABASE_URL=postgres://postgres:admin@localhost:5432/go_template?sslmode=disable make migrate-up

# Generate type-safe database code
make sqlc-generate
```

**Note:** By default, migrations run automatically when the application starts (`DB_AUTO_MIGRATE=true`). You can disable this by setting `DB_AUTO_MIGRATE=false` in your `.env` file.

The API will be available at `http://localhost:8080`

## 🔄 Development Workflows Explained

### Understanding Docker vs Docker Compose

**Docker** - Single container operations:

```bash
docker run postgres:latest    # Run one container
docker build -t myapp .       # Build one image
```

**Docker Compose** - Multi-container applications:

```bash
docker compose up            # Start all services
docker compose up psql_bp    # Start only database service
docker compose down          # Stop all services
```

### When to Use Each Workflow

| Workflow        | Database         | Application             | Use Case                               | Commands                                      |
| --------------- | ---------------- | ----------------------- | -------------------------------------- | --------------------------------------------- |
| **Full Docker** | Docker Container | Docker Container        | Production testing, complete isolation | `make docker-run`                             |
| **Hybrid**      | Docker Container | Local (with hot reload) | Active development, debugging          | `docker compose up psql_bp -d` + `make watch` |
| **Local**       | Docker Container | Local                   | Development without hot reload         | `docker compose up psql_bp -d` + `make run`   |

### Recommended Workflow for Daily Development

1. **Start database once (leave it running):**

   ```bash
   docker compose up psql_bp -d
   ```

2. **Develop with hot reload:**

   ```bash
   make watch
   ```

3. **When done for the day:**
   ```bash
   docker compose down
   ```

### Quick Commands Reference

```bash
# Check what's running
docker compose ps

# View database logs
docker compose logs psql_bp

# Stop only database
docker compose stop psql_bp

# Restart database
docker compose restart psql_bp

# Full cleanup (removes containers and volumes)
docker compose down -v
```

## 📋 Available Commands

### Development Commands

| Command      | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| `make run`   | Run the application                                          |
| `make build` | Build the application binary                                 |
| `make watch` | Live reload during development (auto-installs air if needed) |
| `make test`  | Run all tests                                                |
| `make itest` | Run integration tests only                                   |
| `make clean` | Remove build artifacts                                       |
| `make all`   | Build and test                                               |

### Database Commands

| Command                                   | Description                       |
| ----------------------------------------- | --------------------------------- |
| `make migrate-up`                         | Run database migrations           |
| `make migrate-down`                       | Rollback last migration           |
| `make migrate-status`                     | Check migration status            |
| `make migrate-create name=migration_name` | Create new migration              |
| `make sqlc-generate`                      | Generate Go code from SQL queries |

### Docker Commands

| Command                        | Description                                   |
| ------------------------------ | --------------------------------------------- |
| `make docker-run`              | Start both application and database in Docker |
| `make docker-down`             | Stop all Docker containers                    |
| `docker compose up psql_bp -d` | Start only database container                 |
| `docker compose down`          | Stop all services and remove containers       |

### Tool Installation

| Command              | Description                  |
| -------------------- | ---------------------------- |
| `make install-tools` | Install sqlc and goose tools |

## 🔧 Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and adjust values as needed:

```bash
# Application
APP_NAME=go-template
APP_ENV=development
APP_DEBUG=true
PORT=8080

# Database
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=go_template
BLUEPRINT_DB_USERNAME=postgres
BLUEPRINT_DB_PASSWORD=password
DB_AUTO_MIGRATE=true

# File Upload
UPLOAD_MAX_FILE_SIZE=10485760  # 10MB
UPLOAD_PATH=uploads
BASE_URL=http://localhost:8080

# JWT (for future authentication)
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRES_IN=24h
```

## 📡 API Endpoints

### Health Check

- `GET /api/v1/health` - API health status with database connectivity

### User Management

- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - List users (with pagination)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### File Management

- `POST /api/v1/files/upload` - Upload files
- `GET /api/v1/files` - List all files (with pagination)
- `GET /api/v1/files/my` - List current user's files
- `GET /api/v1/files/:id` - Get file metadata
- `PUT /api/v1/files/:id` - Update file metadata
- `DELETE /api/v1/files/:id` - Delete file
- `GET /api/v1/files/:id/download` - Download file
- `GET /files/:filename` - Serve file directly

For detailed API documentation, see [docs/api/README.md](docs/api/README.md)

## 🔒 Security Features

- **Rate Limiting**: 100 requests per minute per IP
- **Input Validation**: Comprehensive request validation with go-playground/validator
- **File Upload Security**: File type validation, size limits, secure storage
- **CORS**: Configurable cross-origin resource sharing
- **SQL Injection Protection**: Type-safe queries with SQLC
- **Request Logging**: All requests logged with unique request IDs

## 🧪 Testing

```bash
# Run all tests
make test

# Run integration tests only
make itest
```

Integration tests use Testcontainers to spin up real PostgreSQL instances for testing.

## 📦 Tech Stack

- **Language**: Go 1.23.4
- **Framework**: Echo v4
- **Database**: PostgreSQL with pgx driver
- **Query Builder**: SQLC for type-safe SQL
- **Migrations**: Goose
- **Logging**: Zap
- **Validation**: go-playground/validator
- **Testing**: Testcontainers
- **Live Reload**: Air

## 🚀 Production Deployment

The project includes Docker support for easy deployment:

```bash
# Using Docker Compose
docker-compose up --build

# Or build and deploy manually
make build
./main
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🎯 What's Included

This template provides:

- ✅ RESTful API with user and file management
- ✅ Clean, layered architecture
- ✅ Database migrations and type-safe queries
- ✅ Comprehensive middleware (logging, CORS, rate limiting)
- ✅ File upload with validation
- ✅ Structured logging with request tracing
- ✅ Integration testing setup
- ✅ Development tools (live reload, testing)
- ✅ Docker containerization
- ✅ Production-ready configuration
- ✅ API documentation

Perfect for jumpstarting your next Go REST API project! 🎉
