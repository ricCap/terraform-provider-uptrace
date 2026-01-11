variable "uptrace_endpoint" {
  description = "Uptrace API endpoint"
  type        = string
  default     = "http://localhost:14318/internal/v1"
}

variable "uptrace_token" {
  description = "Uptrace API user token"
  type        = string
  sensitive   = true
}

variable "uptrace_project_id" {
  description = "Uptrace project ID"
  type        = number
}
