package docs

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/charmixer/golang-api-template/app"
	"github.com/charmixer/oas/api"

	"github.com/charmixer/golang-api-template/endpoint"

	"github.com/rs/zerolog/log"
)

type GetDocsRequest struct{}
type GetDocsEndpoint struct {
	endpoint.Endpoint
}

func (ep *GetDocsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("request").Start(r.Context(), fmt.Sprintf("%s handler", r.URL.Path))
	defer span.End()

	url := fmt.Sprintf("http://%s:%d/docs/openapi?format=json", app.Env.Domain, app.Env.Port)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error().Err(err)
		panic(err)
	}

	// request.Header.Set("X-Request-Id", ctx.Value("req_id").(string))

	prop := propagation.TraceContext{}
	prop.Inject(ctx, propagation.HeaderCarrier(request.Header))

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

	/*  w.Write([]byte(fmt.Sprintf(`
	<!doctype html> <!-- Important: must specify -->
	<html>
	<head>
	  <meta charset="utf-8">
	  <meta name="viewport" content="width=device-width, minimum-scale=1, initial-scale=1, user-scalable=yes">
	  <link href="https://fonts.googleapis.com/css2?family=Open+Sans:wght@300;600&family=Roboto+Mono&display=swap" rel="stylesheet">
	</head>
	<body>

	  <rapi-doc
	    id="rapidoc-container"
	    theme = "dark"
	    layout = "row"
	    render-style = "read"
	    show-header = "false"
	    allow-try = "false"
	    allow-server-selection = "false"
	    allow-authentication="false"

	    regular-font="Open Sans"
	    mono-font="Roboto Mono"
	  > </rapi-doc>

	  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
	  <script>
	    document.addEventListener('DOMContentLoaded', (event) => {

	      let docEl = document.getElementById("rapidoc-container");

	      let objSpec = JSON.parse(` + "`%s`" + `);
	      docEl.loadSpec(objSpec);
	    })
	  </script>

	</body>
	</html>
		`, spec)))*/

	w.Write([]byte(fmt.Sprintf(`
<!doctype html> <!-- Important: must specify -->
<html>
<head>
	<meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 charecters -->
</head>
<body>

  <div id="redoc-container"></div>

  <script src="https://cdn.jsdelivr.net/npm/redoc@latest/bundles/redoc.standalone.js"> </script>
	<script>
    document.addEventListener('DOMContentLoaded', (event) => {
      let options = {
        hideDownloadButton: true
  		};

  		Redoc.init(JSON.parse(`+"`%s`"+`), options, document.getElementById('redoc-container'))
  	})
	</script>

</body>
</html>
	`, spec)))
}

func NewGetDocsEndpoint() endpoint.EndpointHandler {
	ep := GetDocsEndpoint{}

	ep.Setup(
		endpoint.WithSpecification(api.Path{
			Summary:     "OpenAPI documentation",
			Description: ``,
			Tags:        OPENAPI_TAGS,

			Request: api.Request{
				Description: ``,
				Schema:      GetDocsRequest{},
			},

			Responses: []api.Response{{
				Description: `OpenAPI documentation rendered in HTML`,
				Code:        200,
				ContentType: []string{"text/html"},
				//Schema: GetHealthResponse{},
			}},
		}),
	)

	// Must be pointer to allow ServeHTTP method to be used with *Endpoint
	return &ep
}
