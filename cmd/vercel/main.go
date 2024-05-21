package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"guz.one/config"
	"guz.one/internals"
)

type VercelConfig struct {
	OutputDirectory string `json:"outputDirectory"`
}

func main() {
	configPath := flag.String("c", "./vercel.json", "the path to the vercel.json file")
	staticDir := flag.String("s", "./static", "the directory to copy static files from")
	port := flag.Int("p", 8080, "the port to run the server")

	configFile, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Unable to read vercel.json file due to:\n%s", err)
	}

	var c VercelConfig
	err = json.Unmarshal(configFile, &c)
	if err != nil {
		log.Fatalf("Unable to parse vercel.json file due to:\n%s", err)
	}

	w := internals.StaticWriter{
		DistDir:   &c.OutputDirectory,
		StaticDir: staticDir,
		Pages:     config.ROUTES,
		Context:   context.Background(),
		Logger:    *log.Default(),
	}

	log.Print("Writing static files")
	err = w.WriteAll()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Starting server")
	mux := http.NewServeMux()

	config.APIROUTES(mux)
	mux.Handle("/", http.FileServer(http.Dir(c.OutputDirectory)))

	log.Printf("Running server at port: %v", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", *port), mux)
	if err != nil {
		log.Fatalf("Server crashed due to:\n%s", err)
	}
}
