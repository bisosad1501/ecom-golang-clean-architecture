# Development Guide

This guide helps developers set up and contribute to the E-commerce API project.

## 🛠️ Development Setup

### Prerequisites

- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **PostgreSQL 15+**: [Download PostgreSQL](https://www.postgresql.org/download/)
- **Redis 7+**: [Download Redis](https://redis.io/download)
- **Git**: [Download Git](https://git-scm.com/downloads)
- **Docker** (optional): [Download Docker](https://www.docker.com/get-started)

### Local Development

1. **Clone the repository**
```bash
git clone https://github.com/bisosad1501/ecom-golang-clean-architecture.git
cd ecom-golang-clean-architecture
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment**
```bash
cp .env.example .env
# Edit .env with your local configuration
```

4. **Start dependencies with Docker**
```bash
docker-compose up -d postgres redis
```

5. **Run the application**
```bash
go run cmd/api/main.go
```

6. **Verify setup**
```bash
curl http://localhost:8080/health
```

## 🏗️ Project Structure

```
├── cmd/                    # Application entry points
│   └── api/               # Main API application
├── internal/              # Private application code
│   ├── domain/           # Domain layer (Clean Architecture)
│   │   ├── entities/     # Business entities
│   │   ├── repositories/ # Repository interfaces
│   │   └── services/     # Domain services
│   ├── usecases/         # Use cases layer (business logic)
│   ├── infrastructure/   # Infrastructure layer
│   │   ├── database/     # Database implementations
│   │   └── config/       # Configuration management
│   └── delivery/         # Delivery layer
│       └── http/         # HTTP handlers, middleware, routes
├── pkg/                  # Public packages (if any)
├── docs/                 # Documentation
│   ├── postman/         # Postman collections
│   ├── API.md           # API documentation
│   ├── DEPLOYMENT.md    # Deployment guide
│   └── DEVELOPMENT.md   # This file
├── scripts/              # Build and utility scripts
├── .github/              # GitHub workflows
├── docker-compose.yml    # Docker composition
├── Dockerfile           # Docker image definition
├── Makefile            # Build automation
├── go.mod              # Go module definition
└── README.md           # Project overview
```

## 🧱 Architecture Principles

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - Contains business entities and rules
   - Independent of external concerns
   - Defines repository interfaces

2. **Use Cases Layer** (`internal/usecases/`)
   - Contains application business logic
   - Orchestrates data flow between entities
   - Implements business use cases

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Implements repository interfaces
   - Handles external dependencies (database, cache, etc.)
   - Contains configuration management

4. **Delivery Layer** (`internal/delivery/`)
   - Handles HTTP requests/responses
   - Contains middleware and routing
   - Converts between HTTP and domain models

### Dependency Rule

Dependencies point inward:
- Domain layer has no dependencies
- Use cases depend only on domain
- Infrastructure depends on domain and use cases
- Delivery depends on use cases

## 🔧 Development Workflow

### Adding New Features

1. **Start with Domain**
   - Define entities in `internal/domain/entities/`
   - Add repository interfaces in `internal/domain/repositories/`
   - Create domain services if needed

2. **Implement Use Cases**
   - Create use case interfaces and implementations in `internal/usecases/`
   - Define request/response models
   - Implement business logic

3. **Add Infrastructure**
   - Implement repository interfaces in `internal/infrastructure/database/`
   - Add database migrations if needed
   - Update configuration if required

4. **Create Delivery Layer**
   - Add HTTP handlers in `internal/delivery/http/handlers/`
   - Update routes in `internal/delivery/http/routes/`
   - Add middleware if needed

5. **Update Documentation**
   - Update API documentation
   - Add Postman collection entries
   - Update README if needed

### Code Style Guidelines

1. **Go Conventions**
   - Follow [Effective Go](https://golang.org/doc/effective_go.html)
   - Use `gofmt` for formatting
   - Use `golint` for linting

2. **Naming Conventions**
   - Use descriptive names
   - Follow Go naming conventions
   - Use consistent terminology across layers

3. **Error Handling**
   - Use domain-specific errors
   - Wrap errors with context
   - Handle errors at appropriate layers

4. **Testing**
   - Write unit tests for business logic
   - Use table-driven tests
   - Mock external dependencies

### Database Migrations

1. **Creating Migrations**
```bash
# Manual migration (using GORM AutoMigrate)
# Add new fields to entities and run the application
```

2. **Migration Best Practices**
   - Always backup before migrations
   - Test migrations on development first
   - Use transactions for complex migrations
   - Document breaking changes

### Environment Configuration

```env
# Development
APP_ENV=development
APP_PORT=8080
LOG_LEVEL=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=postgres
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=development-secret-key
JWT_EXPIRE_HOURS=24
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/usecases/...
```

### Test Structure

```go
func TestUserUseCase_Register(t *testing.T) {
    tests := []struct {
        name    string
        request RegisterRequest
        want    *UserResponse
        wantErr bool
    }{
        {
            name: "valid registration",
            request: RegisterRequest{
                Email:     "test@example.com",
                Password:  "password123",
                FirstName: "John",
                LastName:  "Doe",
            },
            want: &UserResponse{
                Email:     "test@example.com",
                FirstName: "John",
                LastName:  "Doe",
            },
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Mocking

Use interfaces for mocking dependencies:

```go
type MockUserRepository struct {
    users map[uuid.UUID]*entities.User
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
    m.users[user.ID] = user
    return nil
}
```

## 🔍 Debugging

### Logging

```go
// Use structured logging
log.WithFields(log.Fields{
    "user_id": userID,
    "action":  "create_order",
}).Info("Creating new order")
```

### Database Debugging

```bash
# Connect to database
docker exec -it ecom_postgres psql -U postgres -d postgres

# View tables
\dt

# Check specific table
SELECT * FROM users LIMIT 5;
```

### API Debugging

```bash
# Test endpoints with curl
curl -X GET http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -v

# Use Postman collection for comprehensive testing
```

## 📦 Building and Deployment

### Local Build

```bash
# Build binary
go build -o bin/ecom-api cmd/api/main.go

# Run binary
./bin/ecom-api
```

### Docker Build

```bash
# Build Docker image
docker build -t ecom-api .

# Run with Docker
docker run -p 8080:8080 ecom-api
```

### Cross-Platform Build

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/ecom-api-linux cmd/api/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/ecom-api-windows.exe cmd/api/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o bin/ecom-api-macos cmd/api/main.go
```

## 🤝 Contributing

### Pull Request Process

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Follow coding standards
   - Add tests for new functionality
   - Update documentation

4. **Test your changes**
   ```bash
   go test ./...
   go build ./...
   ```

5. **Commit your changes**
   ```bash
   git commit -m "feat: add new feature description"
   ```

6. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request**

### Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

### Code Review Checklist

- [ ] Code follows project conventions
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No breaking changes (or properly documented)
- [ ] Performance impact considered
- [ ] Security implications reviewed

## 🚨 Troubleshooting

### Common Issues

1. **Port already in use**
   ```bash
   # Find process using port 8080
   lsof -i :8080
   
   # Kill process
   kill -9 <PID>
   ```

2. **Database connection failed**
   ```bash
   # Check if PostgreSQL is running
   docker ps | grep postgres
   
   # Check logs
   docker logs ecom_postgres
   ```

3. **Module not found**
   ```bash
   # Clean module cache
   go clean -modcache
   
   # Download dependencies
   go mod download
   ```

### Getting Help

- Check existing [GitHub Issues](https://github.com/bisosad1501/ecom-golang-clean-architecture/issues)
- Create a new issue with detailed description
- Join our community discussions
- Review documentation and examples
