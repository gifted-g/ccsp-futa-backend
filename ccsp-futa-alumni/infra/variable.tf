variable "region" {
  description = "AWS region"
  type        = string
}

variable "name" {
  description = "Base name for resources"
  type        = string
}

variable "image" {
  description = "Container image (ECR, Docker Hub, etc.)"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnets for ALB"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "Private subnets for ECS tasks"
  type        = list(string)
}

variable "container_port" {
  description = "App container port"
  type        = number
  default     = 8080
}

variable "desired_count" {
  description = "ECS service desired task count"
  type        = number
  default     = 1
}

variable "task_cpu" {
  description = "Fargate task CPU (e.g., 256, 512, 1024)"
  type        = number
  default     = 256
}

variable "task_memory" {
  description = "Fargate task Memory (e.g., 512, 1024, 2048)"
  type        = number
  default     = 512
}

variable "environment" {
  description = "Plain environment variables"
  type        = map(string)
  default     = {}
}

variable "secrets_ssm_arn" {
  description = "Environment variables backed by SSM/Secrets Manager ARNs"
  type        = map(string)
  default     = {}
}

variable "assign_public_ip" {
  description = "Assign public IP to tasks"
  type        = bool
  default     = false
}

variable "log_retention_days" {
  description = "CloudWatch Logs retention"
  type        = number
  default     = 14
}

variable "health_check_path" {
  description = "ALB health check path"
  type        = string
  default     = "/health"
}
