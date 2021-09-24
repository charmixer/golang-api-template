package openapi

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/charmixer/oas/api"
	"github.com/charmixer/oas/exporter"
	"github.com/charmixer/golang-api-template/app"

	"github.com/rs/zerolog/log"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Documentation", Description:""},
	}
)

type GetDocsRequest struct {}
func (req GetDocsRequest) Specification() (api.Path) {
	return api.Path{
		Summary: "Serve the OpenAPI specification for this api",
		Description: ``,
		Tags: OPENAPI_TAGS,

		Request: api.Request{
			Description: ``,
			Schema: GetDocsRequest{},
		},

		Responses: []api.Response{{
			Description: `Returns openapi spec in given format`,
			Code: 200,
			ContentType: []string{"application/json", "application/yaml"},
			//Schema: GetHealthResponse{},
		}},
	}
}
func (req GetDocsRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// HTTP client call to /docs/openapi?format=json
	url := fmt.Sprintf("http://%s:%d/docs/openapi?format=json", app.Env.Domain, app.Env.Port)

	ctx := r.Context()

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Added tracing tile client
	res, err := client.Do(request) // http.DefaultClient
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}
	defer res.Body.Close()

	spec, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Error().Msgf("Status not OK, got: '%d'", res.StatusCode)
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf(`
<!doctype html> <!-- Important: must specify -->
<html>
<head>
	<meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 charecters -->
</head>
<body>

	<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
	<script>
	window.onload = function(){
		Redoc.init(JSON.parse(` + "`%s`" + `), {
			scrollYOffset: 50
		}, document.getElementById('redoc-container'))
	}
	</script>

	<div id="redoc-container"></div>
</body>
</html>
	`, spec)))
}

type GetOpenapiRequest struct {
	// Mode string `query:"mode"`
	Format string `query:"format" oas:"Format returned by the endpoint, eg. json"`
	// XOverrideMethodHeader string `header:"x-override-method-header" oas:"My description goes here"`
	// Debug string `cookie:"debug"`
}
func (req GetOpenapiRequest) Specification() (api.Path) {
	return api.Path{
		Summary: "Serve the docs for this api",
		Description: ``,
		Tags: OPENAPI_TAGS,

		Request: api.Request{
			Description: ``,
			//Params: GetOpenapiDocsParams{},
			//Schema: GetHealthRequest{},
		},

		Responses: []api.Response{{
			Description: `Returns openapi spec in yaml format`,
			Code: 200,
			ContentType: []string{api.CONTENT_TYPE_HTML},
			//Schema: GetHealthResponse{},
		}},
	}
}
func (req GetOpenapiRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t := r.URL.Query().Get("format")

	var spec string
	var err error

  if t == "json" {
		spec, err = exporter.ToJson(app.Env.OpenAPI)
		if err != nil {
			log.Error().Err(err)
		}
  } else {
		spec, err = exporter.ToYaml(app.Env.OpenAPI)
		if err != nil {
			log.Error().Err(err)
		}
	}

	w.Write([]byte(spec))
}

// Redoc or RapiDoc.. I like RapiDoc's design and options, but it can't show array in array (issue created)

/*func GetOpenapiDocs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(`
<!doctype html> <!-- Important: must specify -->
<html>
<head>
  <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 charecters -->
  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
	<link href="https://fonts.googleapis.com/css2?family=Open+Sans:ital,wght@0,300;0,400;0,700;1,300;1,400;1,700&display=swap" rel="stylesheet">
</head>
<body>
  <rapi-doc
		spec-url="http://%s:%d/docs/openapi"
		theme = "dark"
		layout = "row"
		render-style = "read"
		show-header = "false"
		allow-try = "false"
		allow-server-selection = "false"
	> </rapi-doc>
</body>
</html>
	`, app.Env.Domain, app.Env.Port)))
}*/
