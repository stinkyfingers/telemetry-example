# Two Options for Tracing
https://docs.datadoghq.com/tracing/setup_overview/open_standards/

## 1. Send Traces to OpenTelemetry Collector and Export w/ Datadog Exporter 
This is more vendor-agnostic, but more work to set up.

Why 2 steps? OpenTelemetry and Datadog TraceID and SpanID conventions differ:
https://docs.datadoghq.com/tracing/connect_logs_and_traces/opentelemetry/

### Deploy OpenTelemetry Collector
https://opentelemetry.io/docs/collector/
Datadog has an exporter available within the Otel Collector: https://docs.datadoghq.com/tracing/setup_overview/open_standards/otel_collector_datadog_exporter/

### Example of distributed tracing using Go otel:
https://levelup.gitconnected.com/how-to-implement-opentelemetry-and-propagate-trace-among-microservices-dfa1a1a14865

### Run Example
(Note: this example does not set up the Otel Collector. Traces are written to file.)
`go run otel_collector/server1/main.go`

`go run otel_collector/server2/main.go`

`curl localhost:5551`

View traces in `traces.txt` files. TraceID and SpanID are dumped to stdout. Note that an http request across services has the same SpanID. 

## 2. Ingest Traces with Datadog Agent:
Vendor bound to Datadog, but easy.

### Install Datadog Agent
https://docs.datadoghq.com/getting_started/agent/
You will need a Datadog API key. Datadog offers a Free 14-day trial with optional Free Plan after that. The Free Plan is sufficient for testing/developing tracing.

### Distributed Tracing with Datadog and Go
https://docs.datadoghq.com/tracing/setup_overview/custom_instrumentation/go#distributed-tracing

### Run Example
`go run ingest_with_datadog/server1/main.go`

`go run ingest_with_datadog/server2/main.go`

`curl localhost:5551`

View traces at https://app.datadoghq.com/apm/traces.


