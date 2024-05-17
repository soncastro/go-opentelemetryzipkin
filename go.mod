module desafiopraticoobservabilidadeopentelemetry

go 1.21

toolchain go1.22.3

require (
	github.com/gorilla/mux v1.8.1
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/trace/zipkin v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
)

require (
	github.com/openzipkin/zipkin-go v0.2.5 // indirect
	go.opentelemetry.io/otel/metric v0.20.0 // indirect
	go.opentelemetry.io/otel/trace v0.20.0 // indirect
)
