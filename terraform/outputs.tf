# Terraform Outputs for BWENG Project
# Bu dosya, Terraform çalıştıktan sonra gösterilecek bilgileri tanımlar

# Namespace bilgisi
output "namespace_name" {
  description = "Oluşturulan namespace adı"
  value       = kubernetes_namespace.bweng.metadata[0].name
}

# Environment bilgisi
output "environment" {
  description = "Kullanılan environment"
  value       = var.environment
}

# Service URL'leri
output "api_gateway_url" {
  description = "API Gateway erişim URL'i"
  value       = "http://localhost:8082"
}

output "user_service_url" {
  description = "User Service erişim URL'i"
  value       = "http://localhost:8080"
}

output "order_service_url" {
  description = "Order Service erişim URL'i"
  value       = "http://localhost:8081"
}

# Pod sayıları
output "total_pods" {
  description = "Toplam pod sayısı"
  value       = var.user_service_replicas + var.order_service_replicas + var.api_gateway_replicas + 1 # +1 for PostgreSQL
}

# Deployment durumu
output "deployment_status" {
  description = "Deployment durumları"
  value = {
    user_service_replicas = var.user_service_replicas
    order_service_replicas = var.order_service_replicas
    api_gateway_replicas = var.api_gateway_replicas
    postgres_replicas = 1
  }
}

# Komutlar
output "useful_commands" {
  description = "Kullanışlı komutlar"
  value = {
    check_pods = "kubectl get pods -n ${var.namespace}"
    check_services = "kubectl get svc -n ${var.namespace}"
    check_logs = "kubectl logs -n ${var.namespace} <pod-name>"
    port_forward_gateway = "kubectl port-forward -n ${var.namespace} svc/api-gateway-service 8082:8082"
    port_forward_user = "kubectl port-forward -n ${var.namespace} svc/user-service 8080:8080"
    port_forward_order = "kubectl port-forward -n ${var.namespace} svc/order-service 8081:8081"
  }
} 