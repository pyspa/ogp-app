package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/achiku/mux"
)

func main() {
	configFile := flag.String("c", "", "configuration file path")
	flag.Parse()

	cfg, err := NewConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	app, err := NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/").HandlerFunc(app.InputPage)
	r.Methods(http.MethodGet).Path("/p/{id}").HandlerFunc(app.OgpPage)
	r.Methods(http.MethodGet).PathPrefix("/image/").Handler(
		http.StripPrefix("/image/", http.FileServer(http.Dir("data"))))
	r.Methods(http.MethodPost).Path("/image").HandlerFunc(app.CreateImage)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.APIServerPort), r); err != nil {
		log.Fatal(err)
	}
}
