variable "uptrace_token" {
  description = "Uptrace API user token"
  type        = string
  sensitive   = true
}

variable "uptrace_project_id" {
  description = "Uptrace project ID"
  type        = number
}
