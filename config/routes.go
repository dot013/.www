package config

import (
	"net/http"

	"guz.one/api"
	"guz.one/internals"
	"guz.one/pages"
)

var ROUTES = []internals.Page{
	{Path: "index.html", Component: pages.Homepage()},
}

func APIROUTES(mux *http.ServeMux) {
	mux.HandleFunc("/api/hello", api.Hello)
}
