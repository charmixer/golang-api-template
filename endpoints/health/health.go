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

type Bla struct {
	Test1 bool
	Test2 string
	Test3 int
	Test4 float32
}

type GetHealthRequest struct {
	Alive bool `header:"alive" json:"alive_json" xml:"alive_xml" oas:"Tells if bla"`
	Ready bool `json:"ready_json" xml:"ready_xml"`
	Test struct{
		Bla1 Bla `validate:"required,gte=1"`
		Bla2 []Bla
		Bla3 [][]Bla
		Bla4 [][][]Bla
		String1 []string
		String2 [][]string
		String3 [][][]string
		MapBla map[string]Bla
		MapString map[string]string
		MapInt map[int]int
	}

}
type GetHealthResponse []struct {
	Alive bool `json:"alive_json" xml:"alive_xml" oas:"desc:Tells if bla"`
	Ready bool `json:"ready_json" xml:"ready_xml"`
	Test struct{
		Bla1 Bla
		Bla2 []Bla
		Bla3 [][]Bla
		Bla4 [][][]Bla
		String1 []string
		String2 [][]string
		String3 [][][]string
		MapBla map[string]Bla
		MapString map[string]string
		MapInt map[int]int
	}
}
/*
type PostHealthRequest []struct {
	Alive bool `query:"alive" json:"alive_json" xml:"alive_xml" oas:"Tells if bla"`
	Ready bool `json:"ready_json" xml:"ready_xml"`
	Test struct{
		Bla1 Bla
		Bla2 []Bla
		Bla3 [][]Bla
		Bla4 [][][]Bla
		String1 []string
		String2 [][]string
		String3 [][][]string
		MapBla map[string]Bla
		MapString map[string]string
		MapInt map[int]int
	}
}
type PostHealthResponse []struct {
	Alive bool `json:"alive_json" xml:"alive_xml" desc:"Tells if bla"`
	Ready bool `json:"ready_json" xml:"ready_xml"`
	Test struct{
		Bla1 Bla
		Bla2 []Bla
		Bla3 [][]Bla
		Bla4 [][][]Bla
		String1 []string
		String2 [][]string
		String3 [][][]string
		MapBla map[string]Bla
		MapString map[string]string
		MapInt map[int]int
	}
}*/

func GetHealthSpec() (api.Path) {
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
			ContentType: []string{"application/json", "application/yaml"},
			Schema: GetHealthResponse{},
		}},
	}
}
func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get /health\n"))
}

/*func PostHealthSpec() (api.Path) {
	return api.Path{
		Summary: "Test 2",
		Description: `Testing 2`,

		Request: api.Request{
			Description: `Testing Request`,
			Schema: PostHealthRequest{},
		},

		Responses: []api.Response{{
			Description: `Testing OK Response`,
			Code: 200,
			Schema: PostHealthResponse{},
		}},
	}
}
func PostHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get /health\n"))
}*/

type PostHealthRequest struct {
	Alive bool `query:"alive" json:"alive_json" xml:"alive_xml" oas:"Tells if bla"`
	Ready bool `json:"ready_json" xml:"ready_xml"`
}
type PostHealthResponse []struct {
	Alive bool `json:"alive_json" xml:"alive_xml" desc:"Tells if bla"`
	Ready bool `json:"ready_json" xml:"ready_xml"`
}

func (req PostHealthRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get /health\n"))
}
func (req PostHealthRequest) Specification() (api.Path) {
	return api.Path{
		Summary: "Test 2",
		Description: `Testing 2`,
		Tags: OPENAPI_TAGS,

		Request: api.Request{
			Description: `Testing Request`,
			Schema: PostHealthRequest{},
		},

		Responses: []api.Response{{
			Description: `Testing OK Response`,
			Code: 200,
			Schema: PostHealthResponse{},
		}},
	}
}



/*
type PostHealth struct {
	RequestBody PostHealthRequest
	// RequestParams ..
	ResponseBody PostHealthResponse
}
func (ep PostHealth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get /health\n"))
}
func (ep PostHealth) Specification() (api.Path) {
	return api.Path{
		Summary: "Test 2",
		Description: `Testing 2`,

		Request: api.Request{
			Description: `Testing Request`,
			Schema: ep.RequestBody,
		},

		Responses: []api.Response{{
			Description: `Testing OK Response`,
			Code: 200,
			Schema: ep.ResponseBody,
		}},
	}
}
*/
