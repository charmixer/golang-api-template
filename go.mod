module github.com/charmixer/golang-api-template

go 1.16

require (
	github.com/charmixer/envconfig v1.4.1
	github.com/charmixer/go-flags v1.6.0
	github.com/charmixer/oas v0.0.0-20210928125611-e644f6b3ed57
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.9.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/hetiansu5/urlquery v1.2.7
	github.com/julienschmidt/httprouter v1.3.0
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.25.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.25.0
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/exporters/jaeger v1.0.1
	go.opentelemetry.io/otel/sdk v1.0.1
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/charmixer/oas => ../oas
