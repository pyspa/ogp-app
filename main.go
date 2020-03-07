package main

import (
	"flag"
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/mux"
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

	r.Methods(http.MethodGet).Path("/").HandlerFunc(app.IndexPage)
	r.Methods(http.MethodGet).Path("/p/{id}").HandlerFunc(app.OgpPage)
	r.Methods(http.MethodGet).PathPrefix("/image/").Handler(
		http.StripPrefix("/image/", http.FileServer(http.Dir("data"))))

	// static asset
	r.Methods(http.MethodGet).PathPrefix("/js/").Handler(
		http.StripPrefix("/js/", http.FileServer(http.Dir(path.Join("client", "dist", "js")))))
	r.Methods(http.MethodGet).PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(http.Dir(path.Join("client", "dist", "css")))))
	r.Methods(http.MethodGet).PathPrefix("/img/").Handler(
		http.StripPrefix("/img/", http.FileServer(http.Dir(path.Join("client", "dist", "img")))))

	// API
	r.Methods(http.MethodPost).Path("/api/image").Handler(
		loggingMiddleware(http.HandlerFunc(app.CreateImage)))

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
