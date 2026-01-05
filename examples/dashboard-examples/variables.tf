variable "uptrace_endpoint" {
  description = "Uptrace API endpoint"
  type        = string
}

variable "uptrace_token" {
  description = "Uptrace API authentication token"
  type        = string
  sensitive   = true
}

variable "uptrace_project_id" {
  description = "Uptrace project ID"
  type        = number
}
