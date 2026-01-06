#!/bin/bash
set -e

# Send test telemetry data to Uptrace cloud for testing monitors
# Uses curl to send OTLP/HTTP data directly

UPTRACE_DSN="${UPTRACE_DSN:-https://GZTjjbol8NiGzNEEFzI4Dg@api.uptrace.dev?grpc=4317}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PAYLOAD_FILE="${SCRIPT_DIR}/.otlp_trace_payload.json"
RESPONSE_FILE="${SCRIPT_DIR}/.otlp_response.txt"

echo "Sending test telemetry to Uptrace cloud..."
echo "DSN: ${UPTRACE_DSN}"

# Current timestamp in nanoseconds
NOW_NANOS=$(($(date +%s) * 1000000000))

# Create OTLP trace payload
cat > "${PAYLOAD_FILE}" << EOF
{
  "resourceSpans": [
    {
      "resource": {
        "attributes": [
          {
            "key": "service.name",
            "value": { "stringValue": "terraform-provider-test" }
          },
          {
            "key": "service.version",
            "value": { "stringValue": "1.0.0" }
          },
          {
            "key": "deployment.environment",
            "value": { "stringValue": "test" }
          }
        ]
      },
      "scopeSpans": [
        {
          "scope": {
            "name": "terraform-provider-test"
          },
          "spans": [
EOF

# Generate 50 test spans with varying durations
for i in {1..50}; do
  # Random duration between 100ms and 3s (in nanoseconds)
  DURATION=$((100000000 + RANDOM % 2900000000))
  START_TIME=$((NOW_NANOS - DURATION - i * 1000000000))
  END_TIME=$((START_TIME + DURATION))

  # Alternate between server and client spans
  if [ $((i % 2)) -eq 0 ]; then
    SPAN_KIND="SPAN_KIND_SERVER"
  else
    SPAN_KIND="SPAN_KIND_CLIENT"
  fi

  cat >> "${PAYLOAD_FILE}" << EOF
            {
              "traceId": "$(printf '%032x' $RANDOM$RANDOM$RANDOM$RANDOM)",
              "spanId": "$(printf '%016x' $RANDOM$RANDOM)",
              "name": "test-operation-$i",
              "kind": "$SPAN_KIND",
              "startTimeUnixNano": "$START_TIME",
              "endTimeUnixNano": "$END_TIME",
              "attributes": [
                {
                  "key": "test.type",
                  "value": { "stringValue": "monitor-validation" }
                },
                {
                  "key": "test.index",
                  "value": { "intValue": "$i" }
                },
                {
                  "key": "http.method",
                  "value": { "stringValue": "GET" }
                },
                {
                  "key": "http.status_code",
                  "value": { "intValue": "200" }
                }
              ],
              "status": {
                "code": "STATUS_CODE_OK"
              }
            }$([ $i -lt 50 ] && echo "," || echo "")
EOF
done

cat >> "${PAYLOAD_FILE}" << EOF
          ]
        }
      ]
    }
  ]
}
EOF

echo "Sending OTLP trace data..."

# Send to Uptrace OTLP HTTP endpoint
HTTP_CODE=$(curl -s -w "%{http_code}" -o "${RESPONSE_FILE}" \
  -X POST 'https://otlp.uptrace.dev/v1/traces' \
  -H 'Content-Type: application/json' \
  -H "uptrace-dsn: ${UPTRACE_DSN}" \
  -d @"${PAYLOAD_FILE}")

if [ "$HTTP_CODE" == "200" ] || [ "$HTTP_CODE" == "202" ]; then
  echo "✅ Successfully sent test telemetry data"
  echo "   - 50 test spans with varying durations (100ms-3s)"
  echo "   - Spans include server and client kinds"
  echo "   - HTTP status: $HTTP_CODE"
  echo ""
  echo "Data should be available for monitor tests in ~30 seconds"
  echo "Run: task test:cloud"
else
  echo "❌ Failed to send telemetry. HTTP status: $HTTP_CODE"
  echo "Response:"
  cat "${RESPONSE_FILE}"
  exit 1
fi

# Cleanup
rm -f "${PAYLOAD_FILE}" "${RESPONSE_FILE}"
