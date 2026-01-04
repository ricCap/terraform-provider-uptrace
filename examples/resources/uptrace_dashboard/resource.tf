resource "uptrace_dashboard" "example" {
  yaml = <<-YAML
    name: Service Overview
    grid_rows:
      - title: Traffic
        items:
          - title: Request Rate
            metrics:
              - http_requests_total as $requests
            query:
              - per_min(sum($requests))
          - title: Error Rate
            metrics:
              - http_errors_total as $errors
            query:
              - per_min(sum($errors))
      - title: Resources
        items:
          - title: CPU Usage
            metrics:
              - system.cpu.utilization as $cpu
            query:
              - avg($cpu)
          - title: Memory Usage
            metrics:
              - system.memory.usage as $mem
            query:
              - avg($mem)
  YAML
}
