# ConfigMap and Secret Terraform Configuration
# Bu dosya, ConfigMap ve Secret kaynaklarını tanımlar

# Main Application ConfigMap
resource "kubernetes_config_map" "bweng_config" {
  metadata {
    name      = "bweng-config"
    namespace = var.namespace
    labels = {
      app = "bweng-microservices"
    }
  }

  data = {
    # Database Configuration
    DB_HOST = "postgres-service"
    DB_PORT = tostring(var.postgres_port)
    DB_USER = var.postgres_user
    DB_NAME_USER = "${var.postgres_db}_user_db"
    DB_NAME_ORDER = "${var.postgres_db}_order_db"
    DB_SSLMODE = "disable"
    
    # Service Configuration
    USER_SERVICE_PORT = tostring(var.user_service_port)
    ORDER_SERVICE_PORT = tostring(var.order_service_port)
    API_GATEWAY_PORT = tostring(var.api_gateway_port)
    GRPC_PORT = "50051"
    
    # Environment
    GIN_MODE = "release"
    LOG_LEVEL = "info"
  }
}

# Database ConfigMap (for PostgreSQL)
resource "kubernetes_config_map" "bweng_database_config" {
  metadata {
    name      = "bweng-database-config"
    namespace = "bweng-database"
    labels = {
      app = "bweng-database"
    }
  }

  data = {
    POSTGRES_DB = var.postgres_db
    POSTGRES_USER = var.postgres_user
    POSTGRES_PASSWORD = var.postgres_password
  }
}

# Database Secret (for PostgreSQL)
resource "kubernetes_secret" "bweng_database_secret" {
  metadata {
    name      = "bweng-database-secret"
    namespace = "bweng-database"
    labels = {
      app = "bweng-database"
    }
  }

  type = "Opaque"

  data = {
    # base64 encoded password
    POSTGRES_PASSWORD = base64encode(var.postgres_password)
  }
}

# Application Secret
resource "kubernetes_secret" "bweng_app_secret" {
  metadata {
    name      = "bweng-app-secret"
    namespace = var.namespace
    labels = {
      app = "bweng-microservices"
    }
  }

  type = "Opaque"

  data = {
    # base64 encoded password
    DB_PASSWORD = base64encode(var.postgres_password)
  }
} 