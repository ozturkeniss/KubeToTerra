# Terraform Variables for BWENG Project
# Bu dosya, projenizde kullanılacak tüm değişkenleri tanımlar

# Namespace değişkeni
variable "namespace" {
  description = "Kubernetes namespace adı"
  type        = string
  default     = "bweng"
}

# Environment değişkeni
variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "dev"
  
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment değeri dev, staging veya prod olmalıdır."
  }
}

# Image repository değişkeni
variable "image_repository" {
  description = "Docker image repository"
  type        = string
  default     = "bweng"
}

# Image tag değişkeni
variable "image_tag" {
  description = "Docker image tag"
  type        = string
  default     = "latest"
}

# Replica sayıları
variable "user_service_replicas" {
  description = "User service replica sayısı"
  type        = number
  default     = 2
}

variable "order_service_replicas" {
  description = "Order service replica sayısı"
  type        = number
  default     = 2
}

variable "api_gateway_replicas" {
  description = "API Gateway replica sayısı"
  type        = number
  default     = 2
}

# Port değişkenleri
variable "user_service_port" {
  description = "User service port"
  type        = number
  default     = 8080
}

variable "order_service_port" {
  description = "Order service port"
  type        = number
  default     = 8081
}

variable "api_gateway_port" {
  description = "API Gateway port"
  type        = number
  default     = 8082
}

variable "postgres_port" {
  description = "PostgreSQL port"
  type        = number
  default     = 5432
}

# Database değişkenleri
variable "postgres_user" {
  description = "PostgreSQL kullanıcı adı"
  type        = string
  default     = "postgres"
  sensitive   = true
}

variable "postgres_password" {
  description = "PostgreSQL şifresi"
  type        = string
  default     = "password"
  sensitive   = true
}

variable "postgres_db" {
  description = "PostgreSQL veritabanı adı"
  type        = string
  default     = "bweng"
} 