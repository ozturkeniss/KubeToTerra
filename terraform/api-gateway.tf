# API Gateway Terraform Configuration
# Bu dosya, API Gateway'in deployment ve service'ini tanımlar

# API Gateway Deployment
resource "kubernetes_deployment" "api_gateway" {
  metadata {
    name      = "api-gateway-deployment"
    namespace = var.namespace
    labels = {
      app       = "bweng-microservices"
      component = "api-gateway"
    }
  }

  spec {
    replicas = var.api_gateway_replicas

    selector {
      match_labels = {
        app       = "bweng-microservices"
        component = "api-gateway"
      }
    }

    template {
      metadata {
        labels = {
          app       = "bweng-microservices"
          component = "api-gateway"
        }
      }

      spec {
        container {
          name  = "api-gateway"
          image = "${var.image_repository}-api-gateway:${var.image_tag}"
          
          image_pull_policy = "Never"

          port {
            name           = "http"
            container_port = var.api_gateway_port
          }

          # Environment Variables
          env {
            name = "GIN_MODE"
            value_from {
              config_map_key_ref {
                name = "bweng-config"
                key  = "GIN_MODE"
              }
            }
          }

          env {
            name  = "USER_SERVICE_URL"
            value = "http://user-service:${var.user_service_port}"
          }

          env {
            name  = "ORDER_SERVICE_URL"
            value = "http://order-service:${var.order_service_port}"
          }

          # Resource Limits
          resources {
            requests = {
              memory = "64Mi"
              cpu    = "50m"
            }
            limits = {
              memory = "128Mi"
              cpu    = "100m"
            }
          }

          # Health Checks
          liveness_probe {
            http_get {
              path = "/health"
              port = var.api_gateway_port
            }
            initial_delay_seconds = 30
            period_seconds        = 10
            timeout_seconds       = 5
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = var.api_gateway_port
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

# API Gateway Service
resource "kubernetes_service" "api_gateway" {
  metadata {
    name      = "api-gateway-service"
    namespace = var.namespace
    labels = {
      app       = "bweng-microservices"
      component = "api-gateway"
    }
  }

  spec {
    selector = {
      app       = "bweng-microservices"
      component = "api-gateway"
    }

    port {
      name        = "http"
      port        = var.api_gateway_port
      target_port = var.api_gateway_port
      protocol    = "TCP"
    }

    type = "ClusterIP"
  }
} 