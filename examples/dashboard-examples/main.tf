terraform {
  required_providers {
    uptrace = {
      source = "registry.terraform.io/riccap/uptrace"
    }
  }
}

provider "uptrace" {
  endpoint   = var.uptrace_endpoint
  token      = var.uptrace_token
  project_id = var.uptrace_project_id
}

# ============================================================================
# Example 1: Simple Single-Service Dashboard
# ============================================================================

resource "uptrace_dashboard" "service_overview" {
  name = "Service Overview - Simple"

  yaml = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: Request Rate
            columns:
              - service.name: [api-server]
                span.kind: [server]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Requests/min
          - title: Response Time
            columns:
              - service.name: [api-server]
                span.kind: [server]
            chart:
              - type: line
                metric: span.duration|p95
                legend: P95 Latency
  YAML
}

# ============================================================================
# Example 2: Multi-Service Comparison
# ============================================================================

resource "uptrace_dashboard" "service_comparison" {
  name = "Multi-Service Comparison"

  yaml = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: API Gateway vs Backend
            columns:
              - service.name: [api-gateway, backend-service]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Requests/min
                group_by: [service.name]
          - title: Response Time Comparison
            columns:
              - service.name: [api-gateway, backend-service]
            chart:
              - type: line
                metric: span.duration|p95
                legend: P95 Latency
                group_by: [service.name]
      - row:
          - title: Error Rate by Service
            columns:
              - span.status_code: [error]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Errors/min
                group_by: [service.name]
                color: red
          - title: Success Rate
            columns:
              - span.status_code: [ok]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Success/min
                group_by: [service.name]
                color: green
  YAML
}

# ============================================================================
# Example 3: RED Metrics (Rate, Errors, Duration)
# ============================================================================

resource "uptrace_dashboard" "red_metrics" {
  name = "RED Metrics Dashboard"

  yaml = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: Rate - Requests per Minute
            columns:
              - span.system: [http]
                span.kind: [server]
            chart:
              - type: area
                metric: span.count|per_min
                legend: Total Requests
                color: blue
          - title: Errors - Error Rate
            columns:
              - span.status_code: [error]
            chart:
              - type: area
                metric: span.count|per_min
                legend: Errors
                color: red
      - row:
          - title: Duration - P50 Latency
            columns:
              - span.system: [http]
                span.kind: [server]
            chart:
              - type: line
                metric: span.duration|p50
                legend: P50
                color: green
          - title: Duration - P95 & P99 Latency
            columns:
              - span.system: [http]
                span.kind: [server]
            chart:
              - type: line
                metric: span.duration|p95
                legend: P95
                color: orange
              - type: line
                metric: span.duration|p99
                legend: P99
                color: red
  YAML
}

# ============================================================================
# Example 4: Database Performance Dashboard
# ============================================================================

resource "uptrace_dashboard" "database_performance" {
  name = "Database Performance"

  yaml = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: Query Rate by Database
            columns:
              - db.system: [*]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Queries/min
                group_by: [db.system]
          - title: Query Duration
            columns:
              - db.system: [*]
            chart:
              - type: line
                metric: span.duration|p95
                legend: P95 Query Time
                group_by: [db.system]
      - row:
          - title: Slow Queries (>500ms)
            columns:
              - db.system: [*]
                span.duration: ['>500ms']
            chart:
              - type: line
                metric: span.count|per_min
                legend: Slow Queries
                color: orange
          - title: Database Errors
            columns:
              - db.system: [*]
                span.status_code: [error]
            chart:
              - type: line
                metric: span.count|per_min
                legend: DB Errors
                color: red
  YAML
}

# ============================================================================
# Example 5: HTTP Status Code Breakdown
# ============================================================================

resource "uptrace_dashboard" "http_status_codes" {
  name = "HTTP Status Codes"

  yaml = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: 2xx Success Responses
            columns:
              - http.status_code: ['>=200', '<300']
            chart:
              - type: area
                metric: span.count|per_min
                legend: 2xx/min
                color: green
          - title: 3xx Redirects
            columns:
              - http.status_code: ['>=300', '<400']
            chart:
              - type: area
                metric: span.count|per_min
                legend: 3xx/min
                color: blue
      - row:
          - title: 4xx Client Errors
            columns:
              - http.status_code: ['>=400', '<500']
            chart:
              - type: area
                metric: span.count|per_min
                legend: 4xx/min
                color: orange
              - type: line
                metric: span.count|per_min
                legend: 401 Unauthorized
                where: [http.status_code, '=', '401']
                color: yellow
              - type: line
                metric: span.count|per_min
                legend: 404 Not Found
                where: [http.status_code, '=', '404']
                color: purple
          - title: 5xx Server Errors
            columns:
              - http.status_code: ['>=500']
            chart:
              - type: area
                metric: span.count|per_min
                legend: 5xx/min
                color: red
              - type: line
                metric: span.count|per_min
                legend: 500 Internal Error
                where: [http.status_code, '=', '500']
                color: darkred
  YAML
}

# ============================================================================
# Example 6: External Dependencies
# ============================================================================

resource "uptrace_dashboard" "external_dependencies" {
  name = "External Dependencies"

  yaml = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: External HTTP Calls
            columns:
              - span.kind: [client]
                span.system: [http]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Requests/min
                group_by: [http.host]
          - title: External API Latency
            columns:
              - span.kind: [client]
                span.system: [http]
            chart:
              - type: line
                metric: span.duration|p95
                legend: P95 Latency
                group_by: [http.host]
      - row:
          - title: External API Errors
            columns:
              - span.kind: [client]
                span.system: [http]
                span.status_code: [error]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Errors/min
                group_by: [http.host]
                color: red
          - title: Cache Performance
            columns:
              - span.system: [redis, memcached]
            chart:
              - type: line
                metric: span.duration|p95
                legend: Cache Latency
                group_by: [span.system]
  YAML
}

# ============================================================================
# Example 7: Business Metrics
# ============================================================================

resource "uptrace_dashboard" "business_metrics" {
  name = "Business Metrics"

  yaml = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: User Signups
            columns:
              - span.name: [user.signup]
            chart:
              - type: line
                metric: span.count|per_hour
                legend: Signups/hour
                color: green
          - title: Purchases
            columns:
              - span.name: [order.create]
                span.status_code: [ok]
            chart:
              - type: line
                metric: span.count|per_hour
                legend: Purchases/hour
                color: blue
      - row:
          - title: Failed Payments
            columns:
              - span.name: [payment.process]
                span.status_code: [error]
            chart:
              - type: line
                metric: span.count|per_hour
                legend: Failed Payments
                color: red
          - title: Average Order Value
            columns:
              - span.name: [order.create]
            chart:
              - type: line
                metric: order.amount|avg
                legend: Avg Order $
                color: purple
  YAML
}

# ============================================================================
# Example 8: Microservices Map
# ============================================================================

resource "uptrace_dashboard" "microservices_map" {
  name = "Microservices Health Map"

  yaml = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: Gateway Service
            columns:
              - service.name: [api-gateway]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Traffic
              - type: line
                metric: span.duration|p95
                legend: Latency
          - title: Auth Service
            columns:
              - service.name: [auth-service]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Traffic
              - type: line
                metric: span.duration|p95
                legend: Latency
      - row:
          - title: User Service
            columns:
              - service.name: [user-service]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Traffic
              - type: line
                metric: span.duration|p95
                legend: Latency
          - title: Order Service
            columns:
              - service.name: [order-service]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Traffic
              - type: line
                metric: span.duration|p95
                legend: Latency
      - row:
          - title: Payment Service
            columns:
              - service.name: [payment-service]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Traffic
              - type: line
                metric: span.duration|p95
                legend: Latency
          - title: Notification Service
            columns:
              - service.name: [notification-service]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Traffic
              - type: line
                metric: span.duration|p95
                legend: Latency
  YAML
}
