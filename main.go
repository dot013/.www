package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"www/config"
	"www/internals"
)

var logger = log.Default()

func main() {
	staticDir := flag.String("s", "./static", "the directory to copy static files from")
	port := flag.Int("p", 8080, "the port to run the server")

	mux := http.NewServeMux()

	config.APIROUTES(mux)
	for _, route := range config.ROUTES {
		path := "/" + strings.TrimSuffix(route.Path, ".html")
		if path == "/index" {
			continue
		}
		logger.Printf("Registering page route. page=%s route=%s", route.Path, path)

		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("Handling request. path=%s", r.URL.Path)

			w.Header().Add("Content-Type", "text/html")

			err := route.Component.Render(r.Context(), w)
			if err != nil {
				logger.Fatalf("Unable to render route %s due to %s", route.Path, err)
			}
		})
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			logger.Printf("Handling file server request. path=%s", r.URL.Path)
			http.FileServer(http.Dir(*staticDir)).ServeHTTP(w, r)
			return
		}

		logger.Printf("Handling request. path=%s", r.URL.Path)

		w.Header().Add("Content-Type", "text/html")

		index := slices.IndexFunc(config.ROUTES, func(route internals.Page) bool {
			return route.Path == "index.html"
		})
		indexPage := config.ROUTES[index]

		err := indexPage.Component.Render(r.Context(), w)
		if err != nil {
			log.Fatalf("Unable to render index page due to %s", err)
		}
	})

	logger.Printf("Running server at port: %v", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), mux)
	if err != nil {
		logger.Fatalf("Server crashed due to:\n%s", err)
	}
}
