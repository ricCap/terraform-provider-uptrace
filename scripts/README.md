# Scripts

## send-test-telemetry.sh

Sends test telemetry data (spans) to Uptrace cloud for testing monitors.

**Usage:**
```bash
task test:cloud:send-telemetry
# Wait 30 seconds for data processing
task test:cloud
```

**What it does:**
- Generates 50 test spans with varying durations (100ms-3s)
- Sends data via OTLP/HTTP to Uptrace cloud
- Creates `span.duration` metrics for monitor testing

**Note:** The cloud API test (`TestAccMonitorResource_CloudTrendAggregation`) verifies the `trend_agg_func` field is working. The test may show a "Provider produced inconsistent result" error due to query normalization (cloud API reformats UQL queries), but this does NOT indicate a problem with the `trend_agg_func` feature - the monitor is created successfully and the field works correctly.

## Cleanup

Run `task clean` to remove temporary files and build artifacts.
