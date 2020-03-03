package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/achiku/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	go initProfiler()

	configFile := flag.String("c", "", "configuration file path")
	flag.Parse()

	cfg, err := NewConfig(*configFile)
	if err != nil {
		log.Fatal().Msgf("Failed to create a new configuration: %v", err)
	}
	app, err := NewApp(cfg)
	if err != nil {
		log.Fatal().Msgf("Failed to create a new application: %v", err)
	}

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/").HandlerFunc(app.InputPage)
	r.Methods(http.MethodGet).Path("/p/{id}").HandlerFunc(app.OgpPage)
	r.Methods(http.MethodGet).PathPrefix("/image/").Handler(
		http.StripPrefix("/image/", http.FileServer(http.Dir("data"))))
	r.Methods(http.MethodPost).Path("/image").HandlerFunc(app.CreateImage)

	p := fmt.Sprintf(":%s", cfg.APIServerPort)
	switch isTLS(cfg.BaseURL) {
	case false:
		if err := http.ListenAndServe(p, r); err != nil {
			log.Fatal().Msgf("Failed to run HTTP server without TLS: %v", err)
		}
	case true:
		if err := http.ListenAndServeTLS(p, cfg.ServerCertPath, cfg.ServerKeyPath, r); err != nil {
			log.Fatal().Msgf("Failed to run HTTP server with TLS: %v", err)
		}
	}
}
