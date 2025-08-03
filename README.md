# Go REST API Template

A production-ready, scalable Go REST API template built with modern best practices. This template provides a solid foundation for building RESTful APIs with JWT authentication, user management, file upload capabilities, and comprehensive tooling for rapid development.

## 🚀 Features

- **JWT Authentication**: Complete auth system with access + refresh tokens, password hashing
- **Role-Based Access Control (RBAC)**: Three-tier role system (admin, moderator, user)
- **API Documentation**: Interactive Swagger/OpenAPI documentation with request/response examples
- **Email System**: SMTP email service with beautiful HTML templates for verification and password reset
- **Email Verification**: Secure user email verification with token-based authentication
- **Password Reset**: Complete password reset system with secure tokens and email delivery
- **Clean Architecture**: Layered architecture with clear separation of concerns (handler → service → repository → entity)
- **Echo Framework**: High-performance HTTP router and middleware
- **PostgreSQL Integration**: Raw SQL with pgx driver and SQLC for type-safe queries
- **Database Migrations**: Goose for schema versioning with automatic migration on startup
- **Pagination & Filtering**: Comprehensive pagination system with advanced filtering, search, and sorting capabilities
- **File Management**: Secure file upload with validation and user-linked storage
- **Structured Logging**: Zap logger with request tracing
- **Security Features**: JWT auth, RBAC authorization, bcrypt hashing, rate limiting, CORS, input validation, email verification
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
│   ├── middleware/             # Custom middleware (rate limiting, CORS, logging, JWT auth, RBAC)
│   ├── migration/              # Database migration utilities
│   ├── repository/             # Data access layer with raw SQL
│   ├── router/                 # Route definitions
│   ├── server/                 # Server initialization and configuration
│   ├── service/                # Business logic layer
│   └── utils/                  # Utility functions
├── pkg/
│   ├── jwt/                    # JWT token management utilities
│   ├── pagination/             # Pagination utilities and metadata
│   ├── response/               # Standardized API responses
│   ├── storage/                # File storage utilities
│   └── validator/              # Request validation with custom rules
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
| `make swagger-gen`                        | Generate Swagger API documentation |

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

# JWT Authentication
JWT_ACCESS_SECRET=your-super-secret-access-key-change-this-in-production
JWT_REFRESH_SECRET=your-super-secret-refresh-key-change-this-in-production
JWT_ACCESS_EXPIRES_IN=30m
JWT_REFRESH_EXPIRES_IN=168h

# Database
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=go_template
BLUEPRINT_DB_USERNAME=postgres
BLUEPRINT_DB_PASSWORD=password
BLUEPRINT_DB_SCHEMA=public
DB_AUTO_MIGRATE=true
DATABASE_URL=postgres://postgres:password@localhost:5432/go_template?sslmode=disable

# Server Configuration
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=60s

# File Upload
UPLOAD_MAX_FILE_SIZE=10485760  # 10MB
UPLOAD_PATH=uploads
BASE_URL=http://localhost:8080

# Email Configuration (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=your-email@gmail.com
SMTP_FROM_NAME=Go Template API
```

## 📡 API Endpoints

### Health Check (Public)

- `GET /api/v1/health` - API health status with database connectivity

### Authentication (Public)

- `POST /api/v1/auth/register` - User registration with email/password
- `POST /api/v1/auth/login` - User login with email/password
- `POST /api/v1/auth/refresh` - Refresh access token using refresh token
- `GET /api/v1/auth/verify-email` - Verify email address (supports both GET and POST)
- `POST /api/v1/auth/verify-email` - Verify email address via API
- `POST /api/v1/auth/resend-verification` - Resend email verification
- `POST /api/v1/auth/forgot-password` - Request password reset email
- `GET /api/v1/auth/reset-password` - Validate password reset token (from email links)
- `POST /api/v1/auth/reset-password` - Reset password with token

### Authentication (Protected)

- `GET /api/v1/auth/me` - Get current user profile

### User Management (RBAC Protected)

- `POST /api/v1/users` - Create user (Admin only)
- `GET /api/v1/users` - List all users with pagination and filtering (Moderator+ only)
- `GET /api/v1/users/:id` - Get user by ID (Own profile or Admin)
- `PUT /api/v1/users/:id` - Update user (Own profile or Admin)
- `DELETE /api/v1/users/:id` - Delete user (Admin only)

### File Management (RBAC Protected)

- `POST /api/v1/files/upload` - Upload files (All authenticated users)
- `GET /api/v1/files` - List all files with pagination and filtering (Moderator+ only)
- `GET /api/v1/files/my` - List current user's files with pagination (All authenticated users)
- `GET /api/v1/files/:id` - Get file metadata (All authenticated users)
- `PUT /api/v1/files/:id` - Update file metadata (All authenticated users)
- `DELETE /api/v1/files/:id` - Delete file (Moderator+ only)
- `GET /api/v1/files/:id/download` - Download file (All authenticated users)

### Static Files (Public)

- `GET /files/:filename` - Serve file directly
- `GET /uploads/*` - Static file serving

## 🔍 Pagination & Filtering

The API supports comprehensive pagination, filtering, and search functionality for list endpoints.

### Query Parameters

#### Pagination Parameters
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `sort` - Sort field (varies by endpoint)
- `order` - Sort order: ASC or DESC (default: DESC)
- `search` - Search across multiple fields

#### User Filtering Parameters
- `name` - Filter by name (partial match)
- `email` - Filter by email (partial match)
- `role` - Filter by role (exact match: admin, moderator, user)
- `email_verified` - Filter by email verification status (true/false)
- `created_after` - Filter by creation date (RFC3339 format)
- `created_before` - Filter by creation date (RFC3339 format)

#### File Filtering Parameters
- `file_name` - Filter by file name (partial match)
- `mime_type` - Filter by MIME type (exact match)
- `category` - Filter by category (exact match)
- `uploaded_by` - Filter by uploader user ID
- `created_after` - Filter by upload date (RFC3339 format)
- `created_before` - Filter by upload date (RFC3339 format)

### Example Requests

#### Basic Pagination
```bash
# Get second page with 20 users per page
GET /api/v1/users?page=2&limit=20

# Sort users by name in ascending order
GET /api/v1/users?sort=name&order=ASC
```

#### Filtering Examples
```bash
# Find users with "john" in their name or email
GET /api/v1/users?search=john

# Get only admin users
GET /api/v1/users?role=admin

# Get verified users created after 2024-01-01
GET /api/v1/users?email_verified=true&created_after=2024-01-01T00:00:00Z

# Filter files by MIME type
GET /api/v1/files?mime_type=image/jpeg

# Get files uploaded by specific user
GET /api/v1/files/my?file_name=report&category=documents
```

#### Combined Parameters
```bash
# Complex query: Search for users with "admin" in name/email, 
# verified, created in 2024, sorted by creation date
GET /api/v1/users?search=admin&email_verified=true&created_after=2024-01-01T00:00:00Z&created_before=2024-12-31T23:59:59Z&sort=created_at&order=DESC&page=1&limit=50
```

### Response Format

All paginated endpoints return responses in this format:

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "email_verified": true,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-20T14:20:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total_records": 150,
    "total_pages": 15,
    "has_next": true,
    "has_prev": false
  }
}
```

### Supported Sort Fields

#### Users
- `id` - User ID
- `name` - User name
- `email` - Email address
- `created_at` - Creation timestamp
- `role` - User role

#### Files
- `id` - File ID
- `file_name` - File name
- `file_size` - File size in bytes
- `created_at` - Upload timestamp

### API Documentation

Interactive Swagger/OpenAPI documentation is available at `/swagger/index.html` when the server is running.

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI Spec**: `http://localhost:8080/swagger/doc.json`

The documentation includes:
- Complete endpoint descriptions with RBAC requirements
- Request/response schemas and examples
- Authentication methods (Bearer JWT)
- Interactive API testing interface

To regenerate documentation after making changes:
```bash
make swagger-gen
```

## 📊 Performance & Scalability

### Pagination Architecture
- **Database-Level**: LIMIT/OFFSET pagination at PostgreSQL level
- **Type-Safe Queries**: SQLC-generated code with parameter validation
- **Optimized Counting**: Separate count queries for pagination metadata
- **Index Support**: Optimized for common sort fields (id, created_at)
- **Parameter Limits**: Configurable limits to prevent resource exhaustion

### Future Enhancements
- **Cursor-Based Pagination**: Planned for better performance with large datasets
- **Caching Layer**: Upcoming Redis integration for frequently accessed data
- **Full-Text Search**: Enhanced search capabilities with PostgreSQL FTS

## 🔒 Security Features

- **JWT Authentication**: Access + refresh token system with configurable expiration
- **Role-Based Access Control**: Three-tier role system (admin, moderator, user)
- **Email Verification**: Required email verification for sensitive operations with secure token system
- **Password Reset Security**: Secure token-based password reset with 24-hour expiration
- **Password Security**: Bcrypt hashing with cost 12 (OWASP 2025 recommended)
- **Strong Password Requirements**: 8+ chars, uppercase, lowercase, numbers, special characters
- **Protected Routes**: All user and file endpoints require valid JWT tokens with role-based access
- **Rate Limiting**: 100 requests per minute per IP
- **Input Validation**: Comprehensive request validation with custom password rules
- **File Upload Security**: File type validation, size limits, user-linked uploads
- **CORS**: Configurable cross-origin resource sharing
- **SQL Injection Protection**: Type-safe queries with SQLC
- **Email Security**: Protection against email enumeration attacks
- **Request Logging**: All requests logged with structured format and request IDs

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

- ✅ JWT Authentication with access + refresh tokens
- ✅ Role-Based Access Control (RBAC) with admin, moderator, and user roles
- ✅ Interactive Swagger/OpenAPI documentation
- ✅ Email verification system with secure tokens
- ✅ Password reset functionality with email delivery
- ✅ SMTP email service with beautiful HTML templates
- ✅ RESTful API with user and file management
- ✅ Advanced pagination, filtering, and search capabilities
- ✅ Clean, layered architecture
- ✅ Database migrations and type-safe queries
- ✅ Comprehensive middleware (logging, CORS, rate limiting, JWT auth, RBAC, email verification)
- ✅ Secure file upload with user authentication
- ✅ Password security with bcrypt hashing
- ✅ Structured logging with request tracing
- ✅ Integration testing setup
- ✅ Development tools (live reload, testing)
- ✅ Docker containerization
- ✅ Production-ready configuration
- ✅ API documentation

Perfect for jumpstarting your next Go REST API project! 🎉
