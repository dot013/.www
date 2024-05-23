package config

import (
	"net/http"

	"www/api"
	"www/internals"
	"www/pages"
)

var ROUTES = []internals.Page{
	{Path: "index.html", Component: pages.Homepage()},
}

func APIROUTES(mux *http.ServeMux) {
	mux.HandleFunc("/api/hello", api.Hello)
}
