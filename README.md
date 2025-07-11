# Inventory Manager - Microservices Architecture

A comprehensive microservices-based inventory management system built with Go, featuring distributed architecture, event-driven communication, and containerized deployment.

## Architecture Overview

The system is designed as a distributed microservices architecture with the following components:

### Core Services

- **AuthService**: Handles user authentication, authorization, and JWT token management
- **ItemService**: Manages product catalog and item metadata with MongoDB storage
- **InventoryService**: Tracks inventory levels, stock movements, and availability
- **FiberGateway**: API Gateway providing unified REST endpoints and request routing

### Infrastructure Components

- **PostgreSQL**: Relational database for AuthService and InventoryService
- **MongoDB**: Document database for ItemService product catalog
- **Redis**: Caching layer for session management and performance optimization
- **Kafka**: Message broker for event-driven communication between services
- **Zookeeper**: Coordination service for Kafka cluster management

## Technology Stack

- **Language**: Go 1.23
- **Framework**: Fiber v2 (API Gateway), Gin (Services)
- **Databases**: PostgreSQL 15, MongoDB 6.0
- **Cache**: Redis 7
- **Message Broker**: Apache Kafka 7.4.0
- **Containerization**: Docker & Docker Compose
- **Orchestration**: Kubernetes with Helm charts
- **Protocol**: gRPC for inter-service communication
- **Authentication**: JWT-based token system

## Project Structure

```
rancher-manager/
├── api/                    # Protocol Buffer definitions
│   ├── proto/
│   │   ├── authservice/    # Auth service protobuf
│   │   └── itemservice/    # Item service protobuf
├── cmd/                    # Application entry points
│   ├── authservice/
│   ├── itemservice/
│   └── inventoryservice/
├── internal/               # Core business logic
│   ├── authservice/        # Authentication service implementation
│   ├── itemservice/        # Item management service implementation
│   └── inventoryservice/   # Inventory management service implementation
├── fibergateway/           # API Gateway implementation
├── docker/                 # Docker configuration files
├── helm/                   # Kubernetes Helm charts
├── kafka/                  # Kafka utilities and configurations
├── scripts/                # Build and deployment scripts
└── docker-compose.yml      # Local development environment
```

## Quick Start

### Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- Make (optional, for build scripts)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd rancher-manager
   ```

2. **Start the development environment**
   ```bash
   docker-compose up -d
   ```

3. **Verify services are running**
   ```bash
   docker-compose ps
   ```

4. **Access the API Gateway**
   - URL: http://localhost:8080
   - Health check: http://localhost:8080/health

### Service Endpoints

| Service | Port | Health Check |
|---------|------|--------------|
| FiberGateway | 8080 | http://localhost:8080/health |
| AuthService | 8001 | http://localhost:8001/health |
| ItemService | 8002 | http://localhost:8002/health |
| InventoryService | 8003 | http://localhost:8003/health |

## Development

### Building Services

```bash
# Build all services
make build

# Build specific service
make build-authservice
make build-itemservice
make build-inventoryservice
make build-fibergateway
```

### Running Tests

```bash
# Run all tests
make test

# Run tests for specific service
make test-authservice
make test-itemservice
make test-inventoryservice
```

### Code Generation

```bash
# Generate protobuf code
make proto

# Generate mock files
make mocks
```

## Deployment

### Docker Compose (Development)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### Kubernetes (Production)

```bash
# Build and push images
./scripts/build-images.sh

# Deploy to Kubernetes
./scripts/deploy-k8s.sh

# Clean up deployment
./scripts/cleanup-k8s.sh
```

## Configuration

### Environment Variables

Each service can be configured using environment variables:

- `DATABASE_URL`: PostgreSQL connection string
- `MONGO_URI`: MongoDB connection string
- `REDIS_HOST`: Redis server hostname
- `REDIS_PORT`: Redis server port
- `KAFKA_BROKERS`: Comma-separated list of Kafka brokers
- `JWT_SECRET`: Secret key for JWT token signing

### Docker Compose Override

For development-specific configurations, use `docker-compose.override.yml`:

```yaml
version: '3.8'
services:
  authservice:
    environment:
      - LOG_LEVEL=debug
      - JWT_SECRET=dev-secret-key
    volumes:
      - ./logs:/app/logs
```

## API Documentation

### Authentication Endpoints

- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Token refresh
- `GET /auth/profile` - Get user profile

### Item Management Endpoints

- `GET /items` - List all items
- `POST /items` - Create new item
- `GET /items/{id}` - Get item details
- `PUT /items/{id}` - Update item
- `DELETE /items/{id}` - Delete item

### Inventory Endpoints

- `GET /inventory` - List inventory levels
- `POST /inventory/stock` - Update stock levels
- `GET /inventory/{itemId}` - Get item inventory
- `POST /inventory/movement` - Record stock movement

## Monitoring and Observability

- **Health Checks**: Each service exposes `/health` endpoint
- **Metrics**: Prometheus metrics available on `/metrics`
- **Logging**: Structured JSON logging with correlation IDs
- **Tracing**: Distributed tracing with request correlation

## Security

- JWT-based authentication
- Role-based access control (RBAC)
- Input validation and sanitization
- Secure communication between services
- Environment-based configuration management

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions, please open an issue in the repository or contact the development team. 
