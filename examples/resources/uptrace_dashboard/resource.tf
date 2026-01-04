resource "uptrace_dashboard" "example" {
  yaml = <<-YAML
    name: Service Overview
    gridRows:
      - items:
          - type: chart
            title: Request Rate
            params:
              metrics:
                - name: http_requests_total
                  alias: requests
              query: rate(requests[5m])
              chartKind: line
              legend:
                show: true
          - type: chart
            title: Error Rate
            params:
              metrics:
                - name: http_errors_total
                  alias: errors
              query: rate(errors[5m])
              chartKind: line
              legend:
                show: true
      - items:
          - type: gauge
            title: CPU Usage
            params:
              metrics:
                - name: system.cpu.utilization
                  alias: cpu
              query: avg(cpu)
          - type: gauge
            title: Memory Usage
            params:
              metrics:
                - name: system.memory.usage
                  alias: mem
              query: avg(mem)
  YAML
}
