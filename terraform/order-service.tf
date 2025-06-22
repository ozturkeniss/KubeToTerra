# Order Service Terraform Configuration
# Bu dosya, Order Service'in deployment ve service'ini tanımlar

# Order Service Deployment
resource "kubernetes_deployment" "order_service" {
  metadata {
    name      = "order-service-deployment"
    namespace = var.namespace
    labels = {
      app       = "bweng-microservices"
      component = "order-service"
    }
  }

  spec {
    replicas = var.order_service_replicas

    selector {
      match_labels = {
        app       = "bweng-microservices"
        component = "order-service"
      }
    }

    template {
      metadata {
        labels = {
          app       = "bweng-microservices"
          component = "order-service"
        }
      }

      spec {
        container {
          name  = "order-service"
          image = "${var.image_repository}-order-service:${var.image_tag}"
          
          image_pull_policy = "Never"

          port {
            name           = "http"
            container_port = var.order_service_port
          }

          # Environment Variables
          env {
            name = "DB_HOST"
            value_from {
              config_map_key_ref {
                name = "bweng-config"
                key  = "DB_HOST"
              }
            }
          }

          env {
            name = "DB_PORT"
            value_from {
              config_map_key_ref {
                name = "bweng-config"
                key  = "DB_PORT"
              }
            }
          }

          env {
            name = "DB_USER"
            value_from {
              config_map_key_ref {
                name = "bweng-config"
                key  = "DB_USER"
              }
            }
          }

          env {
            name = "DB_PASSWORD"
            value_from {
              secret_key_ref {
                name = "bweng-app-secret"
                key  = "DB_PASSWORD"
              }
            }
          }

          env {
            name = "DB_NAME"
            value_from {
              config_map_key_ref {
                name = "bweng-config"
                key  = "DB_NAME_ORDER"
              }
            }
          }

          env {
            name = "DB_SSLMODE"
            value_from {
              config_map_key_ref {
                name = "bweng-config"
                key  = "DB_SSLMODE"
              }
            }
          }

          env {
            name = "GIN_MODE"
            value_from {
              config_map_key_ref {
                name = "bweng-config"
                key  = "GIN_MODE"
              }
            }
          }

          # gRPC Service Configuration
          env {
            name  = "USER_SERVICE_GRPC_HOST"
            value = "user-service"
          }

          env {
            name = "USER_SERVICE_GRPC_PORT"
            value_from {
              config_map_key_ref {
                name = "bweng-config"
                key  = "GRPC_PORT"
              }
            }
          }

          # Resource Limits
          resources {
            requests = {
              memory = "128Mi"
              cpu    = "100m"
            }
            limits = {
              memory = "256Mi"
              cpu    = "200m"
            }
          }

          # Health Checks
          liveness_probe {
            http_get {
              path = "/health"
              port = var.order_service_port
            }
            initial_delay_seconds = 30
            period_seconds        = 10
            timeout_seconds       = 5
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = var.order_service_port
            }
            initial_delay_seconds = 5
            period_seconds        = 5
            timeout_seconds       = 3
          }
        }
      }
    }
  }
}

# Order Service Service
resource "kubernetes_service" "order_service" {
  metadata {
    name      = "order-service"
    namespace = var.namespace
    labels = {
      app       = "bweng-microservices"
      component = "order-service"
    }
  }

  spec {
    selector = {
      app       = "bweng-microservices"
      component = "order-service"
    }

    port {
      name        = "http"
      port        = var.order_service_port
      target_port = var.order_service_port
      protocol    = "TCP"
    }

    type = "ClusterIP"
  }
} 