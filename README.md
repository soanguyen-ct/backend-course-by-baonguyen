# Backend Course: User Authentication & Image Upload API

This project implements a REST API for user authentication and image upload using Go and the Echo framework.

## Features
- User registration and login
- JWT-based authentication
- User profile management
- Image upload for authenticated users
- **Flexible Storage Backend**: Support for in-memory, MongoDB, and PostgreSQL storage
- **Containerized Deployment**: Docker and Docker Compose support
- **Production Ready**: Database with proper indexing and error handling

## Project Structure
The project follows a clean 3-layer architecture:
- **Controller Layer**: Handles HTTP requests and responses
- **UseCase Layer**: Contains business logic
- **Storage Layer**: Handles data persistence with multiple backends:
  - In-memory storage
  - MongoDB storage
  - PostgreSQL storage

## API Endpoints

### Public Endpoints
- `POST /api/public/register`: Register a new user
- `POST /api/public/login`: Authenticate and receive a JWT token

### Private Endpoints (Requires Authentication)
- `GET /api/private/self`: Get authenticated user's profile
- `POST /api/private/change-password`: Change password for authenticated user
- `POST /api/private/upload`: Upload an image for the authenticated user

## Quick Start

### Using Makefile (Recommended)
```shell
# Show all available commands
make help

# Quick start with in-memory storage
make quick-start-inmemory

# Quick start with MongoDB
make quick-start-mongodb

# Quick start with PostgreSQL
make quick-start-postgres

# Test all storage backends
make test-backends

# Docker deployment with MongoDB
make docker-mongodb

# Docker deployment with PostgreSQL
make docker-postgres
```

### Manual Setup

### Using Go directly
```shell
go mod download
go run main.go
```

### Using Docker
```shell
docker build -t user-image-api .
docker run -p 8090:8090 user-image-api
```

### Using Docker Compose

#### With MongoDB
```shell
# Start application with MongoDB backend
docker-compose -f docker-compose.mongodb.yml up

# Access MongoDB Admin UI at http://localhost:8081 (admin/admin)
```

#### With PostgreSQL (Recommended for Production)
```shell
# Start application with PostgreSQL backend
docker-compose up

# Access PostgreSQL at localhost:5432
```

#### With In-Memory Storage (Development)
```shell
# Start application with in-memory backend
docker-compose -f docker-compose.inmemory.yml up
```

## API Usage Examples

### Register a User
```shell
curl --location 'localhost:8090/api/public/register' \
--header 'Content-Type: application/json' \
--data '{
    "username": "testuser",
    "password": "password123",
    "full_name": "Test User",
    "address": "123 Test St"
}'
```

### Login
```shell
curl --location 'localhost:8090/api/public/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "testuser",
    "password": "password123"
}'
```

### Change Password
```shell
curl --location 'localhost:8090/api/private/change-password' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN' \
--header 'Content-Type: application/json' \
--data '{
    "old_password": "password123",
    "new_password": "newpassword456"
}'
```

### Get User Profile
```shell
curl --location 'localhost:8090/api/private/self' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjdC1iYWNrZW5kLWNvdXJzZSIsInN1YiI6InRlc3R1c2VyIiwiZXhwIjoxNzQ4NTI3NTU4fQ.EMxlIy-Ux8CJKWya7PQ9HMZIrFmA0d-ZKI9-RTp87Tw'
```

### Upload Image
```shell
curl --location 'localhost:8090/api/private/upload' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjdC1iYWNrZW5kLWNvdXJzZSIsInN1YiI6InRlc3R1c2VyIiwiZXhwIjoxNzUxNDY3OTkyfQ.8nj0dHvhsQSqaXsxoQNpBJ8FzYbiZOwluIxrAzlVGlE' \
--form 'image=@"/path/to/your/image.jpg"'
```

## Storage Configuration

The application supports three storage backends:

### In-Memory Storage
```bash
export STORAGE_TYPE=inmemory
go run main.go
```

### MongoDB Storage
```bash
export STORAGE_TYPE=mongodb
export MONGO_URI=mongodb://localhost:27017
export MONGO_DB_NAME=ct_backend_course
go run main.go
```

### PostgreSQL Storage (Default)
```bash
export STORAGE_TYPE=postgres
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=ct_backend_course
export POSTGRES_SSLMODE=disable
go run main.go
```

### Environment Variables
- `STORAGE_TYPE`: `inmemory`, `mongodb`, or `postgres` (default: `postgres`)
- `MONGO_URI`: MongoDB connection string (default: `mongodb://localhost:27017`)
- `MONGO_DB_NAME`: MongoDB database name (default: `ct_backend_course`)
- `POSTGRES_HOST`: PostgreSQL host (default: `localhost`)
- `POSTGRES_PORT`: PostgreSQL port (default: `5432`)
- `POSTGRES_USER`: PostgreSQL username (default: `postgres`)
- `POSTGRES_PASSWORD`: PostgreSQL password (default: `postgres`)
- `POSTGRES_DB`: PostgreSQL database name (default: `ct_backend_course`)
- `POSTGRES_SSLMODE`: PostgreSQL SSL mode (default: `disable`)

## PostgreSQL Implementation

### Quick PostgreSQL Setup
1. **Local Development**:
   ```bash
   # Start PostgreSQL with Docker
   docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=ct_backend_course postgres:14
   
   # Set environment and run
   export STORAGE_TYPE=postgres
   go run main.go
   ```

2. **Production with Docker Compose**:
   ```bash
   # Complete stack with PostgreSQL
   docker-compose up
   ```

### PostgreSQL Features
- ✅ User and image storage with GORM ORM
- ✅ Automatic migrations
- ✅ Database connection pooling
- ✅ Proper foreign key relationships
- ✅ Indexes for performance optimization
- ✅ Comprehensive error handling

## MongoDB Implementation

For detailed information about the MongoDB implementation, see [MongoDB Guide](docs/MONGODB_GUIDE.md).

### Quick MongoDB Setup
1. **Local Development**:
   ```bash
   # Start MongoDB with Docker
   docker run -d --name mongodb -p 27017:27017 mongo:6.0
   
   # Set environment and run
   export STORAGE_TYPE=mongodb
   go run main.go
   ```

## Switching Between Storage Backends

The application is designed to easily switch between storage backends. To switch:

1. **In Code**: In `main.go`, uncomment the appropriate storage initialization code.

2. **Environment Variables**: Set the `STORAGE_TYPE` environment variable:
   ```bash
   # For MongoDB
   export STORAGE_TYPE=mongodb
   
   # For PostgreSQL
   export STORAGE_TYPE=postgres
   
   # For In-Memory
   export STORAGE_TYPE=inmemory
   ```

3. **Docker Compose**: Use the appropriate compose file:
   ```bash
   # For MongoDB
   docker-compose -f docker-compose.mongodb.yml up
   
   # For PostgreSQL
   docker-compose up
   
   # For In-Memory
   docker-compose -f docker-compose.inmemory.yml up
   ```

## Original Tasks
TODO #1: Restructure to 3-layers (storage, usecase, controller) ✅  
TODO #2: Write private API to upload image and save image info ✅  
TODO #3: Build a Docker image ✅  
TODO #4: Write unit tests for that private API ❌

## Additional Implementations
✅ **JWT Authentication**: Complete user registration and login system with password change functionality  
✅ **Image Upload & Serving**: File validation and static file serving  
✅ **Multiple Storage Backends**: Support for in-memory, MongoDB, and PostgreSQL storage  
✅ **Docker Compose**: Production-ready containerized deployment  
✅ **ORM Integration**: GORM for PostgreSQL  
✅ **Development Tools**: Makefile and scripts for easy development workflow  
✅ **Documentation**: Complete implementation guide and API documentation

