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
  yaml = <<-YAML
    schema: v2
    name: Service Overview - Simple
    grid_rows:
      - title: Traffic Metrics
        items:
          - title: Request Rate
            metrics:
              - span.count as $requests
            query:
              - per_min(sum($requests))
            where:
              - service.name
              - =
              - api-server
          - title: Response Time P95
            metrics:
              - span.duration as $duration
            query:
              - p95($duration)
            where:
              - service.name
              - =
              - api-server
  YAML
}

# ============================================================================
# Example 2: Multi-Service Comparison
# ============================================================================

resource "uptrace_dashboard" "service_comparison" {
  yaml = <<-YAML
    schema: v2
    name: Multi-Service Comparison
    grid_rows:
      - title: Request Rates
        items:
          - title: API Gateway Requests
            metrics:
              - span.count as $requests
            query:
              - per_min(sum($requests))
            where:
              - service.name
              - =
              - api-gateway
          - title: Backend Requests
            metrics:
              - span.count as $requests
            query:
              - per_min(sum($requests))
            where:
              - service.name
              - =
              - backend-service
      - title: Response Times
        items:
          - title: API Gateway P95 Latency
            metrics:
              - span.duration as $duration
            query:
              - p95($duration)
            where:
              - service.name
              - =
              - api-gateway
          - title: Backend P95 Latency
            metrics:
              - span.duration as $duration
            query:
              - p95($duration)
            where:
              - service.name
              - =
              - backend-service
  YAML
}

# ============================================================================
# Example 3: RED Metrics (Rate, Errors, Duration)
# ============================================================================

resource "uptrace_dashboard" "red_metrics" {
  yaml = <<-YAML
    schema: v2
    name: RED Metrics Dashboard
    grid_rows:
      - title: Rate & Errors
        items:
          - title: Request Rate
            metrics:
              - span.count as $requests
            query:
              - per_min(sum($requests))
            where:
              - span.system
              - =
              - http
          - title: Error Rate
            metrics:
              - span.count as $errors
            query:
              - per_min(sum($errors))
            where:
              - span.status_code
              - =
              - error
      - title: Duration (Latency)
        items:
          - title: P50 Latency
            metrics:
              - span.duration as $duration
            query:
              - p50($duration)
            where:
              - span.system
              - =
              - http
          - title: P95 Latency
            metrics:
              - span.duration as $duration
            query:
              - p95($duration)
            where:
              - span.system
              - =
              - http
          - title: P99 Latency
            metrics:
              - span.duration as $duration
            query:
              - p99($duration)
            where:
              - span.system
              - =
              - http
  YAML
}

# ============================================================================
# Example 4: Database Performance Dashboard
# ============================================================================

resource "uptrace_dashboard" "database_performance" {
  yaml = <<-YAML
    schema: v2
    name: Database Performance
    grid_rows:
      - title: Query Metrics
        items:
          - title: Query Rate
            metrics:
              - span.count as $queries
            query:
              - per_min(sum($queries))
            where:
              - db.system
              - in
              - [postgresql, mysql]
          - title: Query Duration P95
            metrics:
              - span.duration as $duration
            query:
              - p95($duration)
            where:
              - db.system
              - in
              - [postgresql, mysql]
      - title: Issues
        items:
          - title: Slow Queries (>500ms)
            metrics:
              - span.count as $slow
            query:
              - per_min(sum($slow))
            where:
              - span.duration
              - '>'
              - 500ms
          - title: Database Errors
            metrics:
              - span.count as $errors
            query:
              - per_min(sum($errors))
            where:
              - db.system
              - exists
              - true
              - span.status_code
              - =
              - error
  YAML
}

# ============================================================================
# Example 5: HTTP Status Code Breakdown
# ============================================================================

resource "uptrace_dashboard" "http_status_codes" {
  yaml = <<-YAML
    schema: v2
    name: HTTP Status Codes
    grid_rows:
      - title: Success & Redirects
        items:
          - title: 2xx Success
            metrics:
              - span.count as $success
            query:
              - per_min(sum($success))
            where:
              - http.status_code
              - '>='
              - 200
              - http.status_code
              - <
              - 300
          - title: 3xx Redirects
            metrics:
              - span.count as $redirects
            query:
              - per_min(sum($redirects))
            where:
              - http.status_code
              - '>='
              - 300
              - http.status_code
              - <
              - 400
      - title: Client & Server Errors
        items:
          - title: 4xx Client Errors
            metrics:
              - span.count as $client_errors
            query:
              - per_min(sum($client_errors))
            where:
              - http.status_code
              - '>='
              - 400
              - http.status_code
              - <
              - 500
          - title: 5xx Server Errors
            metrics:
              - span.count as $server_errors
            query:
              - per_min(sum($server_errors))
            where:
              - http.status_code
              - '>='
              - 500
  YAML
}

# ============================================================================
# Example 6: External Dependencies
# ============================================================================

resource "uptrace_dashboard" "external_dependencies" {
  yaml = <<-YAML
    schema: v2
    name: External Dependencies
    grid_rows:
      - title: External API Calls
        items:
          - title: Outbound Request Rate
            metrics:
              - span.count as $requests
            query:
              - per_min(sum($requests))
            where:
              - span.kind
              - =
              - client
              - span.system
              - =
              - http
          - title: External API Latency
            metrics:
              - span.duration as $duration
            query:
              - p95($duration)
            where:
              - span.kind
              - =
              - client
              - span.system
              - =
              - http
      - title: Cache & Errors
        items:
          - title: Cache Hit Rate
            metrics:
              - span.count as $cache_ops
            query:
              - per_min(sum($cache_ops))
            where:
              - span.system
              - in
              - [redis, memcached]
          - title: External API Errors
            metrics:
              - span.count as $errors
            query:
              - per_min(sum($errors))
            where:
              - span.kind
              - =
              - client
              - span.status_code
              - =
              - error
  YAML
}

# ============================================================================
# Example 7: Business Metrics
# ============================================================================

resource "uptrace_dashboard" "business_metrics" {
  yaml = <<-YAML
    schema: v2
    name: Business Metrics
    grid_rows:
      - title: User Activity
        items:
          - title: User Signups
            metrics:
              - span.count as $signups
            query:
              - per_hour(sum($signups))
            where:
              - span.name
              - =
              - user.signup
          - title: Successful Purchases
            metrics:
              - span.count as $purchases
            query:
              - per_hour(sum($purchases))
            where:
              - span.name
              - =
              - order.create
              - span.status_code
              - =
              - ok
      - title: Revenue & Issues
        items:
          - title: Failed Payments
            metrics:
              - span.count as $failures
            query:
              - per_hour(sum($failures))
            where:
              - span.name
              - =
              - payment.process
              - span.status_code
              - =
              - error
          - title: Active Sessions
            metrics:
              - span.count as $sessions
            query:
              - sum($sessions)
            where:
              - span.name
              - =
              - session.active
  YAML
}

# ============================================================================
# Example 8: Microservices Health Map
# ============================================================================

resource "uptrace_dashboard" "microservices_map" {
  yaml = <<-YAML
    schema: v2
    name: Microservices Health Map
    grid_rows:
      - title: Core Services
        items:
          - title: API Gateway
            metrics:
              - span.count as $requests
              - span.duration as $duration
            query:
              - per_min(sum($requests))
              - p95($duration)
            where:
              - service.name
              - =
              - api-gateway
          - title: Auth Service
            metrics:
              - span.count as $requests
              - span.duration as $duration
            query:
              - per_min(sum($requests))
              - p95($duration)
            where:
              - service.name
              - =
              - auth-service
      - title: Business Logic
        items:
          - title: User Service
            metrics:
              - span.count as $requests
              - span.duration as $duration
            query:
              - per_min(sum($requests))
              - p95($duration)
            where:
              - service.name
              - =
              - user-service
          - title: Order Service
            metrics:
              - span.count as $requests
              - span.duration as $duration
            query:
              - per_min(sum($requests))
              - p95($duration)
            where:
              - service.name
              - =
              - order-service
      - title: Infrastructure
        items:
          - title: Payment Service
            metrics:
              - span.count as $requests
              - span.duration as $duration
            query:
              - per_min(sum($requests))
              - p95($duration)
            where:
              - service.name
              - =
              - payment-service
          - title: Notification Service
            metrics:
              - span.count as $requests
              - span.duration as $duration
            query:
              - per_min(sum($requests))
              - p95($duration)
            where:
              - service.name
              - =
              - notification-service
  YAML
}
