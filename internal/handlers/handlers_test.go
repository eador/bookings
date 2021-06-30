package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "get", []postData{}, http.StatusOK},
	{"about", "/about", "get", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "get", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "get", []postData{}, http.StatusOK},
	{"contact", "/contact", "get", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "get", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "get", []postData{}, http.StatusOK},

	{"post-search-avail-json", "/search-availability-json", "post", []postData{
		{key: "start", value: "2020-01-02"},
		{key: "end", value: "2020-02-02"},
	}, http.StatusOK},
	{"post-search-avail", "/search-availability", "post", []postData{
		{key: "start", value: "2020-01-02"},
		{key: "end", value: "2020-02-02"},
	}, http.StatusOK},
	{"post-make-reservation", "/make-reservation", "post", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "j.smith@smith.com"},
		{key: "phone", value: "(555) 867-5309"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "get" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
