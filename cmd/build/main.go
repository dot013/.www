package main

import (
	"context"
	"flag"
	"log"

	"guz.one/internals"
)

func main() {
	dir := flag.String("d", "./dist", "the directory to write the files")
	staticDir := flag.String("s", "./static", "the directory to copy static files from")

	w := internals.StaticWriter{
		DistDir:   dir,
		StaticDir: staticDir,
		Pages:     internals.ROUTES,
		Context:   context.Background(),
		Logger:    *log.Default(),
	}

	err := w.WriteAll()
	if err != nil {
		log.Fatal(err)
	}
}
