# BWENG - Microservices Architecture with Kubernetes & Terraform

A comprehensive microservices platform built with Go, featuring user management and order processing capabilities. This project demonstrates modern cloud-native development practices with containerization, orchestration, and infrastructure as code.

## 🏗️ Architecture Overview


![ChatGPT Image Jun 22, 2025, 09_16_15 PM](https://github.com/user-attachments/assets/49674514-4f53-4143-a8cd-cbae8ece6700)


## 🚀 Technology Stack

### Backend Services
- **Language**: Go 1.24.4
- **Framework**: Gin (HTTP server)
- **Database**: PostgreSQL 15 with GORM
- **Communication**: gRPC for inter-service communication
- **Documentation**: Swagger/OpenAPI

### Infrastructure & DevOps
- **Containerization**: Docker
- **Orchestration**: Kubernetes (Minikube)
- **Infrastructure as Code**: Terraform
- **Service Discovery**: Kubernetes Services
- **Load Balancing**: Kubernetes LoadBalancer

### Development Tools
- **API Documentation**: Swagger UI
- **Protocol Buffers**: gRPC definitions
- **Configuration Management**: Kubernetes ConfigMaps & Secrets
- **Health Checks**: Liveness & Readiness probes

## 📁 Project Structure

![Uploading ChatGPT Image Jun 22, 2025, 09_10_45 PM.png…]()



## 🎯 Core Features

### User Management Service
- **User Registration & Authentication**
- **Profile Management** (CRUD operations)
- **Data Validation** with Gin binding
- **Database Integration** with GORM
- **gRPC Communication** for inter-service calls

### Order Processing Service
- **Order Creation & Management**
- **Status Tracking** (Pending, Confirmed, Shipped, Delivered, Cancelled)
- **User Integration** via gRPC calls
- **Business Logic** for order processing
- **Data Consistency** with transactions

### API Gateway
- **Request Routing** to appropriate microservices
- **Load Balancing** across service instances
- **CORS Support** for cross-origin requests
- **Health Monitoring** and service discovery
- **Request Logging** and monitoring

### Infrastructure Management
- **Kubernetes Orchestration** with proper resource allocation
- **Terraform Automation** for infrastructure provisioning
- **Persistent Storage** for database data
- **Service Discovery** and networking
- **Health Checks** and auto-scaling

## 🛠️ Quick Start

### Prerequisites
- Go 1.24.4+
- Docker & Docker Compose
- Kubernetes (Minikube)
- Terraform 1.0+

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd bweng
   ```

2. **Start Kubernetes cluster**
   ```bash
   minikube start --driver=docker
   ```

3. **Build and deploy with Terraform**
   ```bash
   cd terraform
   terraform init
   terraform plan
   terraform apply
   ```

4. **Access the services**
   ```bash
   # API Gateway
   kubectl port-forward -n bweng svc/api-gateway-service 8082:8082
   
   # User Service
   kubectl port-forward -n bweng svc/user-service 8080:8080
   
   # Order Service
   kubectl port-forward -n bweng svc/order-service 8081:8081
   ```

### API Endpoints

#### User Service (Port 8080)
```
GET    /health                    # Health check
GET    /api/v1/users             # List all users
POST   /api/v1/users             # Create user
GET    /api/v1/users/:id         # Get user by ID
PUT    /api/v1/users/:id         # Update user
DELETE /api/v1/users/:id         # Delete user
```

#### Order Service (Port 8081)
```
GET    /health                    # Health check
GET    /api/v1/orders            # List all orders
POST   /api/v1/orders            # Create order
GET    /api/v1/orders/:id        # Get order by ID
PUT    /api/v1/orders/:id/status # Update order status
DELETE /api/v1/orders/:id        # Delete order
```

#### API Gateway (Port 8082)
```
GET    /health                    # Gateway health check
GET    /services                  # List registered services
GET    /api/v1/users/*           # Proxy to user service
GET    /api/v1/orders/*          # Proxy to order service
```

## 🔧 Configuration

### Environment Variables
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `GIN_MODE`: Gin framework mode
- `USER_SERVICE_GRPC_HOST`: User service gRPC host
- `USER_SERVICE_GRPC_PORT`: User service gRPC port

### Kubernetes Configuration
- **Namespaces**: `bweng` (main), `bweng-database` (PostgreSQL)
- **Replicas**: 2 for each microservice (configurable)
- **Resource Limits**: CPU and memory constraints defined
- **Health Checks**: Liveness and readiness probes configured

## 📊 Monitoring & Observability

### Health Checks
- **Liveness Probes**: Ensure services are running
- **Readiness Probes**: Ensure services are ready to serve traffic
- **Database Connectivity**: PostgreSQL health monitoring

### Logging
- **Structured Logging**: JSON format for better parsing
- **Request Tracing**: Gateway-level request logging
- **Error Tracking**: Comprehensive error handling and logging

### Metrics
- **Service Metrics**: Response times and throughput
- **Resource Utilization**: CPU and memory usage
- **Database Performance**: Query performance monitoring

## 🔒 Security Features

- **Input Validation**: Request validation with Gin binding
- **SQL Injection Prevention**: GORM ORM with parameterized queries
- **Secret Management**: Kubernetes secrets for sensitive data
- **Network Security**: Kubernetes network policies
- **CORS Configuration**: Cross-origin resource sharing setup

## 🚀 Deployment Strategies

### Development Environment
- **Local Kubernetes**: Minikube for development
- **Hot Reloading**: Development mode with auto-restart
- **Debugging**: Integrated debugging support

### Production Environment
- **High Availability**: Multiple replicas for fault tolerance
- **Auto-scaling**: Horizontal Pod Autoscaler (HPA)
- **Load Balancing**: Kubernetes service load balancing
- **Persistent Storage**: Reliable data storage with PVCs

## 📈 Scalability

### Horizontal Scaling
- **Stateless Services**: Easy horizontal scaling
- **Database Scaling**: Read replicas and connection pooling
- **Load Distribution**: Kubernetes service mesh capabilities

### Performance Optimization
- **Connection Pooling**: Database connection management
- **Caching**: Redis integration ready
- **CDN Integration**: Static content delivery optimization

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go coding standards
- Write comprehensive tests
- Update documentation
- Use conventional commit messages

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the documentation in the `docs/` folder
- Review the Kubernetes and Terraform configurations

## 🔮 Roadmap

### Planned Features
- [ ] Authentication & Authorization (JWT)
- [ ] Redis caching layer
- [ ] Message queue integration (RabbitMQ/Kafka)
- [ ] Prometheus metrics collection
- [ ] Grafana dashboards
- [ ] CI/CD pipeline with GitHub Actions
- [ ] Multi-region deployment
- [ ] Blue-green deployment strategy

### Technical Improvements
- [ ] Service mesh implementation (Istio)
- [ ] Advanced monitoring with Jaeger tracing
- [ ] Database migrations with Flyway
- [ ] Automated testing with TestContainers
- [ ] Security scanning with Trivy
- [ ] Performance benchmarking suite

---

**Built with ❤️ using Go, Kubernetes, and Terraform** 
