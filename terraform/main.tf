# Terraform Configuration for BWENG Microservices Project
# Bu dosya, Kubernetes cluster'ınızda çalışan tüm servisleri tanımlar

terraform {
  required_version = ">= 1.0"
  
  # Kubernetes provider'ını kullanacağız
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.25"
    }
  }
}

# Kubernetes provider konfigürasyonu
provider "kubernetes" {
  # Minikube context'ini kullanıyoruz
  config_path = "~/.kube/config"
}

# Namespace oluştur
resource "kubernetes_namespace" "bweng" {
  metadata {
    name = var.namespace
  }
} 