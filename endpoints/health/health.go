package health

import (
	"net/http"
	"github.com/charmixer/oas/api"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Health", Description:""},
	}
)
/*
type Request interface{}
type Response interface {
	Response() interface{}
}
type Endpoint interface {
	http.Handler
	Specification() (api.Api)
	Request()
	Response()
	Validate(Request) error
	Validate(Response) error
}

// https://golang.org/doc/effective_go#embedding
type HealthEndpoint struct {
	Endpoint
	Request
	Response
}
func (ep *HealthEndpoint) Request() GetHealthRequest {

	return nil
}
func (ep *HealthEndpoint) Response() IResponse {

	return nil
}*/

type GetHealthRequest struct {}
type GetHealthResponse []struct {
	Alive bool `json:"alive_json" xml:"alive_xml" desc:"Tells if bla"`
	Ready bool `json:"ready_json" xml:"ready_xml"`
}

func (req GetHealthRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get /health\n"))
}
func (req GetHealthRequest) Specification() (api.Path) {
	return api.Path{
		Summary: "Test 2",
		Description: `Testing 2`,
		Tags: OPENAPI_TAGS,

		Request: api.Request{
			Description: `Testing Request`,
			Schema: GetHealthRequest{},
		},

		Responses: []api.Response{{
			Description: `Testing OK Response`,
			Code: 200,
			Schema: GetHealthResponse{},
		}},
	}
}
