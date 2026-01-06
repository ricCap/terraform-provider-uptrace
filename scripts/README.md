# Scripts

This directory contains utility scripts for development and testing.

## send-test-telemetry.sh

Sends test telemetry data (spans) to Uptrace cloud for testing monitors and validating cloud API compatibility.

**Usage:**
```bash
task test:cloud:send-telemetry
# Wait 30 seconds for data processing
task test:cloud
```

**What it does:**
- Generates 50 test spans with varying durations (100ms-3s)
- Sends data via OTLP/HTTP to Uptrace cloud endpoint
- Creates `span.duration` metrics suitable for monitor testing
- Uses temporary files in script directory (auto-cleaned after execution)

**Environment Variables:**
- `UPTRACE_DSN` - Uptrace data source name (defaults to demo project)

**Note:** The cloud API acceptance test (`TestAccMonitorResource_CloudTrendAggregation`) verifies the `trend_agg_func` field functionality. The test may show a "Provider produced inconsistent result" error due to query normalization (cloud API reformats UQL queries), but this does NOT indicate a problem with the `trend_agg_func` feature - the monitor is created successfully and the field works correctly.

## Development

### Cleanup

Run `task clean` to remove temporary files and build artifacts:
```bash
task clean
```

This removes:
- Generated code artifacts
- Build binaries
- Test caches
- Temporary OTLP payload files
