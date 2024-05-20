package main

import (
	"context"
	"log"
	"net/http"

	"guz.one/api"
	"guz.one/pages"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/hello", api.Hello)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := pages.Homepage().Render(context.Background(), w)
		_ = err
	})

	log.Fatal(http.ListenAndServe(":5432", mux))
}
