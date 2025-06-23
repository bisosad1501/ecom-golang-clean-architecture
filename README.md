# E-commerce Golang Clean Architecture

A modern e-commerce system built with Go following Clean Architecture principles.

## 🏗️ Architecture

This project follows Clean Architecture with the following layers:

```
├── cmd/                    # Application entry points
│   └── api/               # Main API application
├── internal/
│   ├── domain/            # Domain layer (entities, repositories, services)
│   │   ├── entities/      # Business entities
│   │   ├── repositories/  # Repository interfaces
│   │   └── services/      # Domain services
│   ├── usecases/          # Use cases layer (business logic)
│   ├── infrastructure/    # Infrastructure layer
│   │   ├── database/      # Database implementations
│   │   └── config/        # Configuration
│   └── delivery/          # Delivery layer
│       └── http/          # HTTP handlers, middleware, routes
├── pkg/                   # Shared packages
├── docs/                  # Documentation
├── scripts/               # Build and deployment scripts
└── migrations/            # Database migrations
```

## ✨ Features

- **User Management**: Registration, Authentication, Profile Management
- **Product Management**: CRUD operations, Categories, Search functionality
- **Shopping Cart**: Add/Remove items, Update quantities
- **Order Management**: Create orders, Track status, Order history
- **Payment Processing**: Multiple payment methods support
- **Inventory Management**: Stock tracking and management
- **Admin Panel**: User and product management
- **API Security**: JWT authentication, Role-based access control

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **HTTP Framework**: Gin
- **ORM**: GORM
- **Authentication**: JWT
- **Validation**: Go Playground Validator
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **CI/CD**: GitHub Actions

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

### Quick Start with Docker

1. Clone the repository
```bash
git clone https://github.com/bisosad1501/ecom-golang-clean-architecture.git
cd ecom-golang-clean-architecture
```

2. Start with Docker Compose
```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Manual Installation

1. Clone the repository
```bash
git clone https://github.com/bisosad1501/ecom-golang-clean-architecture.git
cd ecom-golang-clean-architecture
```

2. Install dependencies
```bash
go mod tidy
```

3. Set up environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Start PostgreSQL and Redis
```bash
docker-compose up -d postgres redis
```

5. Run the application
```bash
make run
# or
go run cmd/api/main.go
```

## 📚 API Documentation

### Available Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login

#### User Management
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile
- `POST /api/v1/users/change-password` - Change password

#### Admin Endpoints
- `GET /api/v1/admin/users` - List all users
- `POST /api/v1/admin/users/:id/activate` - Activate user
- `POST /api/v1/admin/users/:id/deactivate` - Deactivate user

### Example Usage

#### Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### Get user profile (requires authentication)
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Build the application
make build

# Clean build artifacts
make clean
```

## 🐳 Docker

### Build Docker image
```bash
make docker-build
```

### Run with Docker Compose
```bash
docker-compose up -d
```

### Services included:
- **API**: Main application (port 8080)
- **PostgreSQL**: Database (port 5432)
- **Redis**: Cache (port 6379)
- **pgAdmin**: Database management (port 5050) - optional

## 🔧 Configuration

Key environment variables:

```env
# Application
APP_NAME=ecom-api
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=postgres

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRE_HOURS=24
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Clean Architecture by Robert C. Martin
- Go community for excellent packages and tools
- Contributors and maintainers