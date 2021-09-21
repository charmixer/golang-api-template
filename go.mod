module github.com/charmixer/golang-api-template

go 1.16

require (
	github.com/charmixer/envconfig v1.4.1
	github.com/charmixer/go-flags v1.6.0
	github.com/charmixer/oas v0.0.0-20200807123054-f614693b91d8
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/justinas/alice v1.2.0
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.24.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.23.0
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC3
	go.opentelemetry.io/otel/sdk v1.0.0-RC3
	go.opentelemetry.io/otel/trace v1.0.0-RC3
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/charmixer/oas => /Users/laroed/git/charmixer/oas
