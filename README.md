# Golang Service

A RESTful API service built with Go, featuring Azure Entra ID authentication, PostgreSQL database, Docker containerization, and Kubernetes deployment.

## Features

- 🚀 **RESTful API** with Gin framework
- 🔐 **Azure Entra ID Authentication** for secure access
- 🗄️ **PostgreSQL Database** with GORM ORM
- 🐳 **Docker Support** for containerization
- ☸️ **Kubernetes Ready** with deployment manifests
- 📝 **OpenAPI Specification** with Swagger documentation
- 🏥 **Health Check Endpoint** for monitoring
- 🔧 **Environment-based Configuration**

## Project Structure

```
├── cmd/server/              # Application entry point
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── database/          # Database setup and migrations
│   ├── handlers/          # HTTP request handlers
│   ├── middleware/        # HTTP middleware
│   └── models/           # Data models
├── api/                   # OpenAPI specification
├── deployments/           # Deployment configurations
│   ├── docker/           # Docker files
│   └── kubernetes/       # Kubernetes manifests
├── scripts/              # Build and deployment scripts
└── docs/                 # Documentation
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (for local development without Docker)
- kubectl (for Kubernetes deployment)

### Local Development

1. **Setup environment:**
   ```bash
   make dev-setup
   ```

2. **Update configuration:**
   Edit `.env` file with your Azure Entra ID credentials:
   ```bash
   AZURE_TENANT_ID=your-actual-tenant-id
   AZURE_CLIENT_ID=your-actual-client-id
   JWT_SECRET=your-super-secret-jwt-key
   ```

3. **Run with Docker Compose:**
   ```bash
   make docker-compose-up
   ```

4. **Or run locally:**
   ```bash
   # Start PostgreSQL locally
   # Update DATABASE_URL in .env to point to local instance
   make run
   ```

### API Documentation

Once running, access the Swagger documentation at:
- Local: http://localhost:8080/swagger/index.html

### Health Check

Check service health:
```bash
curl http://localhost:8080/health
```

## API Endpoints

### Health
- `GET /health` - Health check (no authentication required)

### Users (requires authentication)
- `GET /api/v1/users` - List all users
- `POST /api/v1/users` - Create new user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

## Authentication

The service uses Azure Entra ID for authentication. Include the JWT token in the Authorization header:

```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
     http://localhost:8080/api/v1/users
```

## Docker Usage

### Build Image
```bash
make docker-build
```

### Run Container
```bash
make docker-run
```

### Docker Compose (Recommended for local development)
```bash
# Start services
make docker-compose-up

# View logs
make docker-compose-logs

# Stop services
make docker-compose-down
```

## Kubernetes Deployment

### Prerequisites
- Kubernetes cluster
- kubectl configured
- Docker image pushed to registry

### Deploy to Kubernetes

1. **Update image references** in `deployments/kubernetes/deployment.yaml`

2. **Update secrets** in `deployments/kubernetes/secret.yaml` with base64 encoded values:
   ```bash
   echo -n "your-actual-tenant-id" | base64
   ```

3. **Deploy:**
   ```bash
   make k8s-deploy
   ```

4. **Check status:**
   ```bash
   make k8s-logs
   kubectl get pods -n golang-service
   ```

### Configuration

The service uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Runtime environment | `development` |
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | See .env.example |
| `AZURE_TENANT_ID` | Azure Entra ID Tenant ID | Required |
| `AZURE_CLIENT_ID` | Azure Entra ID Client ID | Required |
| `JWT_SECRET` | JWT signing secret | Required |

## Development

### Available Commands

```bash
make help                    # Show all available commands
make build                   # Build the application
make run                     # Run the application
make test                    # Run tests
make fmt                     # Format code
make lint                    # Run linter
make check                   # Run all checks (fmt, lint, test)
make swagger                 # Generate Swagger documentation
```

### Adding New Endpoints

1. Define models in `internal/models/`
2. Add handlers in `internal/handlers/`
3. Register routes in `cmd/server/main.go`
4. Update OpenAPI spec in `api/openapi.yaml`

## Production Considerations

### Security
- Use proper Azure Entra ID token validation
- Store secrets in Kubernetes secrets or cloud key vaults
- Enable TLS/SSL termination at load balancer
- Implement rate limiting
- Use read-only filesystem in containers

### Monitoring
- Health check endpoint: `/health`
- Add metrics endpoint (Prometheus)
- Implement structured logging
- Add distributed tracing

### Scalability
- Horizontal pod autoscaling configured
- Database connection pooling
- Stateless design for easy scaling

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make check`
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For questions or issues, please create an issue in the repository.
