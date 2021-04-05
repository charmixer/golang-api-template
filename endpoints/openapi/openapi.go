package openapi

import (
	"fmt"
	"net/http"
	"github.com/charmixer/oas/api"

	"github.com/charmixer/golang-api-template/app"
)

type GetOpenapiDocsParams struct {
	Mode string `query:"mode"`
	Test []map[string]string `query:"test" header:"x-test" path:"test"`
	XOverrideMethodHeader string `header:"x-override-method-header" oas:"My description goes here"`
	Debug string `cookie:"debug"`
}
type GetOpenapiDocsRequest struct {
	Mode string
	Test []string
	Test2 []struct{
		A string
	}
}

func GetOpenapiSpec() (api.Path) {
	return api.Path{
		Summary: "Serve the OpenAPI specification for this api",
		Description: ``,

		Request: api.Request{
			Description: ``,
			//Schema: GetHealthRequest{},
		},

		Responses: []api.Response{{
			Description: `Returns openapi spec in yaml format`,
			Code: 200,
			ContentType: api.CONTENT_TYPE_TEXT,
			//Schema: GetHealthResponse{},
		}},
	}
}
func GetOpenapi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(app.Env.OpenAPI))
}

func GetOpenapiDocsSpec() (api.Path) {
	return api.Path{
		Summary: "Serve the docs for this api",
		Description: ``,

		Request: api.Request{
			Description: ``,
			Params: GetOpenapiDocsParams{},
			//Schema: GetHealthRequest{},
		},

		Responses: []api.Response{{
			Description: `Returns openapi spec in yaml format`,
			Code: 200,
			ContentType: api.CONTENT_TYPE_HTML,
			//Schema: GetHealthResponse{},
		}},
	}
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
		spec-url="http://%s:%d/docs/openapi.yaml"
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
func GetOpenapiDocs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(`
<!doctype html> <!-- Important: must specify -->
<html>
<head>
  <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 charecters -->
</head>
<body>
  <redoc
		spec-url="http://%s:%d/docs/openapi.yaml"
		hide-loading="false"
		path-in-middle-panel="false"
	></redoc>
	<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
</body>
</html>
	`, app.Env.Domain, app.Env.Port)))
}
