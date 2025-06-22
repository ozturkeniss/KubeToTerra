# PostgreSQL Terraform Configuration
# Bu dosya, PostgreSQL'in deployment, service ve storage kaynaklarını tanımlar

# PostgreSQL Namespace
resource "kubernetes_namespace" "bweng_database" {
  metadata {
    name = "bweng-database"
    labels = {
      app = "bweng-database"
    }
  }
}

# PostgreSQL Persistent Volume
resource "kubernetes_persistent_volume" "bweng_postgres_pv" {
  metadata {
    name = "bweng-postgres-pv"
    labels = {
      app = "bweng-database"
    }
  }

  spec {
    capacity = {
      storage = "10Gi"
    }
    access_modes = ["ReadWriteOnce"]
    persistent_volume_source {
      host_path {
        path = "/data/bweng-postgres"
      }
    }
    storage_class_name = "local-storage"
  }
}

# PostgreSQL Persistent Volume Claim
resource "kubernetes_persistent_volume_claim" "bweng_postgres_pvc" {
  metadata {
    name      = "bweng-postgres-pvc"
    namespace = var.namespace
    labels = {
      app = "bweng-database"
    }
  }

  spec {
    access_modes = ["ReadWriteOnce"]
    storage_class_name = "local-storage"
    resources {
      requests = {
        storage = "10Gi"
      }
    }
  }
}

# PostgreSQL Deployment
resource "kubernetes_deployment" "postgres" {
  metadata {
    name      = "postgres-deployment"
    namespace = var.namespace
    labels = {
      app       = "bweng-database"
      component = "postgres"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app       = "bweng-database"
        component = "postgres"
      }
    }

    template {
      metadata {
        labels = {
          app       = "bweng-database"
          component = "postgres"
        }
      }

      spec {
        container {
          name  = "postgres"
          image = "postgres:15-alpine"

          port {
            container_port = var.postgres_port
          }

          # Environment Variables
          env {
            name = "POSTGRES_DB"
            value_from {
              config_map_key_ref {
                name = "bweng-database-config"
                key  = "POSTGRES_DB"
              }
            }
          }

          env {
            name = "POSTGRES_USER"
            value_from {
              config_map_key_ref {
                name = "bweng-database-config"
                key  = "POSTGRES_USER"
              }
            }
          }

          env {
            name = "POSTGRES_PASSWORD"
            value_from {
              secret_key_ref {
                name = "bweng-database-secret"
                key  = "POSTGRES_PASSWORD"
              }
            }
          }

          # Volume Mounts
          volume_mount {
            name       = "postgres-storage"
            mount_path = "/var/lib/postgresql/data"
          }

          volume_mount {
            name       = "init-script"
            mount_path = "/docker-entrypoint-initdb.d"
          }

          # Resource Limits
          resources {
            requests = {
              memory = "256Mi"
              cpu    = "250m"
            }
            limits = {
              memory = "512Mi"
              cpu    = "500m"
            }
          }

          # Health Checks
          liveness_probe {
            exec {
              command = ["pg_isready", "-U", "postgres"]
            }
            initial_delay_seconds = 30
            period_seconds        = 10
          }

          readiness_probe {
            exec {
              command = ["pg_isready", "-U", "postgres"]
            }
            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }

        # Volumes
        volume {
          name = "postgres-storage"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.bweng_postgres_pvc.metadata[0].name
          }
        }

        volume {
          name = "init-script"
          config_map {
            name = "postgres-init-script"
          }
        }
      }
    }
  }
}

# PostgreSQL Service
resource "kubernetes_service" "postgres" {
  metadata {
    name      = "postgres-service"
    namespace = var.namespace
    labels = {
      app       = "bweng-database"
      component = "postgres"
    }
  }

  spec {
    selector = {
      app       = "bweng-database"
      component = "postgres"
    }

    port {
      name        = "postgres"
      port        = var.postgres_port
      target_port = var.postgres_port
      protocol    = "TCP"
    }

    type = "ClusterIP"
  }
} 