# Getting started

After cloning repo, run:

```
go mod init <new repo eg. github.com/charmixer/golang-api-template>
go mod tidy
go run main.go serve
```

# What the template gives

## Configuration setup

Reading configuration in following order

1. Read from yaml file by setting `CFG_PATH=/path/to/conf.yaml`
2. Read from environment `PATH_TO_STRUCT_FIELD=override-value`
3. Read from flags `go run main.go serve -f override-value`
4. If none of above use `default:"value"` tag from configuration struct

## Metrics

Middleware for prometheus has been added. Access metrics from `/metrics`

Includes defaults and stuff like total requests. Customize in `middleware/metrics.go`

## Documentation

Using the structs to setup your endpoints will allow for automatic generation of openapi spec.

Documentation can be found at `/docs` and spec at `/docs/openapi.yaml`
If you get an error try setting `-d localhost`

# TODO / Ideas

- [ ] Setup changelog generator (https://github.com/charmixer/auto-changelog-action)
- [x] Setup tracing with OpenTelemetry and Jaeger example
- [ ] Setup input / output validation with (https://github.com/go-playground/validator) - must be easy to remove
- [ ] Show validation in docs (OpenAPI spec)
- [ ] README.md update with guides
- [ ] Setup some sort of error handling
- [ ] Utils for parsing request body and query
