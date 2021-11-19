# Getting started

After cloning repo, run:

```
find . -type f \( -name "*.go" \) -exec sed -i '' 's/charmixer\/golang-api-template/your-repo-name/g' {} +;
go mod init <new repo eg. github.com/username/your-repo-name>
go mod tidy
go run main.go serve
```

To build with supported ldflags use:
```
go build -ldflags="-s -w -X main.version=1.0.0 -X main.commit=qwerty -X main.date=20210101 -X main.tag=v1.0.0 -X main.name=golang-template-api -X main.environment=production" .
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

# Done / TODOs

- [x] Automated documentation of endpoints (OpenAPI)
- [x] Metrics endpoint with default collectors (Prometheus)
- [x] Struct-based application config
  - [x] Config from yaml
  - [x] Override config with environment
  - [x] Override environment with flags
  - [x] Default configuration values
- [x] Setup tracing with OpenTelemetry and Jaeger example
- [x] Setup input / output validation with (https://github.com/go-playground/validator) - must be easy to remove/disable
  - [x] Input validation
  - [x] Output validation
- [x] Setup some sort of error handling ([rfc7807](https://datatracker.ietf.org/doc/html/rfc7807))
- [x] Utils for parsing request body and query
- [x] Show validation in docs (OpenAPI spec)
- [x] Health checks with uptime and external deps setup (https://datatracker.ietf.org/doc/html/draft-inadarei-api-health-check)
- [x] CI pipeline with fmt, vet, staticchecks, build & test
- [x] CD pipeline using goreleaser - triggered by tag push
- [x] Publish to Github Packages (ghcr.io) with GoReleaser - disable by adding skip flag to the docker section in: `.goreleaser.yml` (see GoReleaser docs)
- [x] Setup changelog generator (https://github.com/charmixer/auto-changelog-action) - currently creates pull requests
  - [ ] Should somehow create the changelog before release tag is created, so it gets baked in
- [ ] README.md update with guides
- [ ] HTTP Client with easy tracing propagation
- [ ] Infrastructure reference stack
