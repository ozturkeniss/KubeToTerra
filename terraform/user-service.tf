# User Service Terraform Configuration
# Bu dosya, User Service'in deployment ve service'ini tanımlar

# User Service Deployment
resource "kubernetes_deployment" "user_service" {
  metadata {
    name      = "user-service-deployment"
    namespace = var.namespace
    labels = {
      app       = "bweng-microservices"
      component = "user-service"
    }
  }

  spec {
    replicas = var.user_service_replicas

    selector {
      match_labels = {
        app       = "bweng-microservices"
        component = "user-service"
      }
    }

    template {
      metadata {
        labels = {
          app       = "bweng-microservices"
          component = "user-service"
        }
      }

      spec {
        container {
          name  = "user-service"
          image = "${var.image_repository}-user-service:${var.image_tag}"
          
          image_pull_policy = "Never"

          port {
            name           = "http"
            container_port = var.user_service_port
          }

          port {
            name           = "grpc"
            container_port = 50051
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
                key  = "DB_NAME_USER"
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
              port = var.user_service_port
            }
            initial_delay_seconds = 30
            period_seconds        = 10
            timeout_seconds       = 5
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = var.user_service_port
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

# User Service Service
resource "kubernetes_service" "user_service" {
  metadata {
    name      = "user-service"
    namespace = var.namespace
    labels = {
      app       = "bweng-microservices"
      component = "user-service"
    }
  }

  spec {
    selector = {
      app       = "bweng-microservices"
      component = "user-service"
    }

    port {
      name        = "http"
      port        = var.user_service_port
      target_port = var.user_service_port
      protocol    = "TCP"
    }

    port {
      name        = "grpc"
      port        = 50051
      target_port = 50051
      protocol    = "TCP"
    }

    type = "ClusterIP"
  }
} 