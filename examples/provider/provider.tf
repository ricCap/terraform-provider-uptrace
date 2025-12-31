terraform {
  required_providers {
    uptrace = {
      source = "riccap/uptrace"
    }
  }
}

provider "uptrace" {
  endpoint   = "https://api2.uptrace.dev"
  token      = var.uptrace_token
  project_id = var.uptrace_project_id
}

variable "uptrace_token" {
  description = "Uptrace API token"
  type        = string
  sensitive   = true
}

variable "uptrace_project_id" {
  description = "Uptrace project ID"
  type        = number
}
