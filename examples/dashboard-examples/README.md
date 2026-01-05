# Dashboard Examples

This directory contains 8 comprehensive dashboard examples demonstrating different visualization patterns and use cases for the Uptrace Terraform provider.

## Dashboard Gallery

### 1. **Simple Single-Service Dashboard**
Basic dashboard for monitoring one service.

**Use Case**: Individual service health check
**Metrics**: Request rate, response time
**Best For**: Getting started, small services

### 2. **Multi-Service Comparison**
Compare metrics across multiple services side-by-side.

**Use Case**: Service dependency analysis
**Metrics**: Traffic, latency, errors by service
**Best For**: Microservices architecture

### 3. **RED Metrics (Rate, Errors, Duration)**
Industry-standard RED metrics dashboard.

**Use Case**: SRE golden signals monitoring
**Metrics**: Request rate, error rate, P50/P95/P99 latency
**Best For**: Production monitoring, SLO tracking

### 4. **Database Performance**
Monitor database queries and performance.

**Use Case**: Database health and optimization
**Metrics**: Query rate, duration, slow queries, errors
**Best For**: Database-heavy applications

### 5. **HTTP Status Code Breakdown**
Detailed breakdown of HTTP response codes.

**Use Case**: API health and debugging
**Metrics**: 2xx/3xx/4xx/5xx counts, specific codes (401, 404, 500)
**Best For**: API services, troubleshooting

### 6. **External Dependencies**
Monitor third-party API calls and cache performance.

**Use Case**: Dependency health tracking
**Metrics**: External API latency/errors, cache performance
**Best For**: Services with external dependencies

### 7. **Business Metrics**
Track business KPIs alongside technical metrics.

**Use Case**: Business intelligence and product analytics
**Metrics**: Signups, purchases, failed payments, order value
**Best For**: Product teams, business stakeholders

### 8. **Microservices Map**
Health map showing all microservices in one view.

**Use Case**: System-wide health overview
**Metrics**: Traffic and latency for each service
**Best For**: Large microservices architectures

## Quick Start

```bash
# Copy example config
cp terraform.tfvars.example terraform.tfvars

# Edit with your credentials
vim terraform.tfvars

# Apply
terraform init
terraform apply
```

## Dashboard YAML Structure

All dashboards use the Uptrace YAML format:

```yaml
table:
  - schema_version: v2
    row:
      - title: Chart Title
        columns:
          - service.name: [your-service]  # Filters
        chart:
          - type: line                     # Chart type
            metric: span.count|per_min     # Metric to plot
            legend: Label                  # Legend text
            group_by: [attribute]          # Optional grouping
            color: blue                    # Optional color
```

## Chart Types

- **`line`**: Time series line chart (default)
- **`area`**: Filled area chart
- **`bar`**: Bar chart

## Common Metrics

### Request Metrics
- `span.count|per_min` - Requests per minute
- `span.count|per_hour` - Requests per hour
- `span.count|per_sec` - Requests per second

### Duration Metrics
- `span.duration|p50` - 50th percentile (median)
- `span.duration|p95` - 95th percentile
- `span.duration|p99` - 99th percentile
- `span.duration|avg` - Average duration
- `span.duration|max` - Maximum duration

### Custom Metrics
- `attribute.name|avg` - Average of any numeric attribute
- `attribute.name|sum` - Sum of any numeric attribute
- `attribute.name|max` - Maximum of any numeric attribute

## Common Filters

### By Service
```yaml
columns:
  - service.name: [api-server, backend]
```

### By HTTP Method
```yaml
columns:
  - http.method: [GET, POST]
```

### By Status Code
```yaml
columns:
  - span.status_code: [error]
  - http.status_code: ['>=500']
```

### By Database
```yaml
columns:
  - db.system: [postgresql, mysql]
```

### Multiple Conditions (AND)
```yaml
columns:
  - service.name: [api-server]
    span.kind: [server]
    span.status_code: [error]
```

## Grouping

Group metrics by any attribute:

```yaml
chart:
  - type: line
    metric: span.count|per_min
    legend: Traffic
    group_by: [service.name]  # Separate line per service
```

Common group_by attributes:
- `service.name` - By service
- `http.method` - By HTTP method
- `http.status_code` - By status code
- `db.system` - By database type
- `deployment.environment` - By environment

## Colors

Available colors:
- `blue`, `green`, `red`, `orange`, `purple`, `yellow`
- `darkred`, `darkblue`, `darkgreen`

## Advanced Patterns

### Multiple Metrics in One Chart

```yaml
chart:
  - type: line
    metric: span.duration|p50
    legend: P50
    color: green
  - type: line
    metric: span.duration|p95
    legend: P95
    color: orange
  - type: line
    metric: span.duration|p99
    legend: P99
    color: red
```

### Filtered Metrics

```yaml
chart:
  - type: line
    metric: span.count|per_min
    legend: Total Errors
  - type: line
    metric: span.count|per_min
    legend: 404 Errors
    where: [http.status_code, '=', '404']
```

### Multi-Row Dashboards

```yaml
table:
  - schema_version: v2
    row:
      - title: Row 1, Column 1
        # ...
      - title: Row 1, Column 2
        # ...
    row:
      - title: Row 2, Column 1
        # ...
      - title: Row 2, Column 2
        # ...
```

## Customization Tips

### 1. Adjust to Your Service Names

Replace example service names with yours:

```hcl
yaml = <<-YAML
  columns:
    - service.name: [your-actual-service-name]
YAML
```

### 2. Add Environment Filters

```yaml
columns:
  - deployment.environment: [production]
    service.name: [api-server]
```

### 3. Customize Time Aggregation

```yaml
chart:
  - metric: span.count|per_min    # Per minute
  - metric: span.count|per_hour   # Per hour
  - metric: span.count|per_sec    # Per second
```

### 4. Focus on Critical Services

Create separate dashboards for:
- Payment/checkout flows
- Authentication systems
- Data pipelines
- External integrations

## Testing Your Dashboard

Before deploying:

1. **Verify YAML Syntax**
   ```bash
   terraform validate
   ```

2. **Check Indentation**
   - Use 2-space indentation
   - No tabs
   - Consistent nesting

3. **Test Filters in Uptrace UI**
   - Go to Uptrace → Metrics
   - Run your filter queries
   - Verify data exists

4. **Apply and Review**
   ```bash
   terraform apply
   # Then check dashboard in Uptrace UI
   ```

## Troubleshooting

### Dashboard Shows No Data

**Cause**: Filters don't match your data

**Fix**:
1. Check attribute names in Uptrace UI
2. Verify service names are correct
3. Ensure time range has data

### YAML Parse Error

**Cause**: Invalid YAML syntax

**Fix**:
1. Check indentation (must be consistent)
2. Verify no tabs (use spaces only)
3. Ensure proper nesting

### Chart Not Rendering

**Cause**: Invalid metric or aggregation

**Fix**:
1. Verify metric exists: `span.duration`, `span.count`, etc.
2. Check aggregation: `|p95`, `|per_min`, `|avg`
3. Test in Uptrace UI first

## Performance Considerations

### Dashboard Load Time

**Large time ranges + many services = slow loading**

Optimize:
```yaml
# ❌ Avoid: Too broad
columns:
  - span.system: [*]  # All systems

# ✅ Better: Specific
columns:
  - service.name: [api-server, backend]  # Only what you need
```

### Too Many Charts

**More than 10-12 charts = slow dashboard**

Solution:
- Split into multiple dashboards
- Use grouping instead of separate charts
- Focus on critical metrics

## Dashboard Organization

Recommended structure:

```
1. Executive Dashboard
   - High-level KPIs
   - Overall system health
   - Business metrics

2. Service Dashboards (one per service)
   - Service-specific RED metrics
   - Dependencies
   - Errors

3. Infrastructure Dashboard
   - Database performance
   - Cache performance
   - External dependencies

4. Troubleshooting Dashboard
   - Detailed error breakdown
   - Slow queries
   - Status codes
```

## Related Examples

- [Complete Setup](../complete-setup/) - Full monitoring stack
- [Multi-Channel Alerts](../multi-channel-alerts/) - Alert routing patterns
- [Getting Started Guide](../../docs/guides/getting-started.md) - Basic setup

## Support

- [Uptrace Dashboard Documentation](https://uptrace.dev/docs/dashboards/)
- [Provider Documentation](../../docs/)
- [GitHub Issues](https://github.com/riccap/tofu-uptrace-provider/issues)
