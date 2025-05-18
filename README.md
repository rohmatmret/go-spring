# Go Spring

A Go application that emulates Spring Boot's structure and functionality, focusing on features like dependency injection, configuration management, caching, and observability.

## Features

- **Dependency Injection**:
- **Configuration Management**: Using HashiCorp Vault
- **Caching**: Support for memory and Redis with `@Cacheable` equivalent
- **Observability**: Prometheus metrics and OpenTelemetry tracing
- **Database**: PostgreSQL with transaction management
- **API**: RESTful endpoints with proper error handling

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- HashiCorp Vault (for configuration)
- PostgreSQL (included in Docker setup)
- Redis (optional, for distributed caching)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-spring.git
   cd go-spring
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   export VAULT_ADDR="http://vault:8200"
   export VAULT_TOKEN="your-vault-token"
   export VAULT_MOUNT_PATH="secret"
   export VAULT_SECRET_PATH="go-spring/config"
   ```



4. Update the configuration in `config.json`:
   ```json
   {
       "server": {
           "port": 8080
       },
       "database": {
           "url": "postgres://postgres:postgres@localhost:5432/go_spring?sslmode=disable"
       },
       "cache": {
           "type": "memory",
           "redis": {
               "url": "redis://localhost:6379/0"
           }
       }
   }
   ```

## Start Infrastructure
**Start the infrastructure services**
   ```bash
   make up
   ```

## Running the Application

1. Start the application:
   ```bash
   go run main.go
   ```

2. The server will start on port 8080.

## Infrastructure Services

The project includes several infrastructure services managed by Docker Compose:

- **PostgreSQL**: Database
  - Port: 5432
  - User: go_spring
  - Password: go_spring_pass
  - Database: go_spring_db

- **Prometheus**: Metrics collection
  - Port: 9090
  - URL: http://localhost:9090

- **Grafana**: Metrics visualization
  - Port: 3000
  - URL: http://localhost:3000
  - Default credentials: admin/admin

- **InfluxDB**: Long-term metrics storage
  - Port: 8086
  - URL: http://localhost:8086

## Available Make Commands

```bash
make up      # Start all services
make down    # Stop all services
make logs    # View service logs
make ps      # Show running containers
make clean   # Remove all containers and volumes
make restart # Restart all services
make health  # Check service health
```

## API Endpoints

### Health Check
- `GET /health`: Check application health

### Metrics
- `GET /metrics`: Prometheus metrics endpoint

### Users
- `POST /api/users`: Create a new user
- `GET /api/users/{id}`: Get a user by ID
- `GET /api/users?username={username}`: Get a user by username
- `PUT /api/users/{id}`: Update a user

## Observability

### Metrics
The application exposes Prometheus metrics at `/metrics`. Key metrics include:
- HTTP request duration and count
- Cache hit/miss rates
- Service method duration

### Tracing
The application uses OpenTelemetry for distributed tracing. Traces can be viewed in Jaeger.

## Project Structure

```
.
├── config.json           # Application configuration
├── go.mod               # Go module file
├── go.sum               # Go module checksum
├── main.go             # Application entry point
└── internal/
    ├── cache/          # Caching implementation
    ├── config/         # Configuration management
    ├── container/      # Dependency injection
    ├── handler/        # HTTP handlers
    ├── observability/  # Metrics and tracing
    ├── repository/     # Data access layer
    └── service/        # Business logic
```

## Contributing


1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 