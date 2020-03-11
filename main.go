package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// parseCommandLineAndReadConfig parses the command line arguments and read the configuration file
func parseCommandLineAndReadConfig() *Config {
	configFile := flag.String("c", "", "configuration file path")
	bind := flag.String("bind", "", "host:port pair where the server binds (defaults to PORT environment variable)")
	baseURL := flag.String("baseurl", "", "Base URL")
	tls := flag.Bool("tls", false, "enable TLS")
	certificatePath := flag.String("cert", "", "path to server certificate file (PEM format)")
	privKeyPath := flag.String("key", "", "path to key file that corresponds to the server certificate (PEM format)")
	flag.Parse()

	cfg, err := NewConfig(*configFile)
	if err != nil {
		log.Fatal().Msgf("Failed to create a new configuration: %v", err)
	}

	if *bind != "" {
		cfg.APIServerPort = ""
		cfg.APIServerBind = *bind
	}

	if *baseURL != "" {
		cfg.BaseURL = *baseURL
	}
	if *tls {
		cfg.TLS = true
	}
	if *certificatePath != "" {
		cfg.ServerCertPath = *certificatePath
	}
	if *privKeyPath != "" {
		cfg.ServerKeyPath = *privKeyPath
	}

	port := os.Getenv("PORT")
	if port != "" {
		cfg.APIServerPort = ""
		cfg.APIServerBind = fmt.Sprintf(":%s", port)
	}

	// for backwards compatibility (TBR)
	if cfg.APIServerPort != "" {
		if cfg.APIServerBind == "" {
			logger.Warn().Msg("api_server_port option is deprecated; use api_server_bind instead.")
			cfg.APIServerBind = fmt.Sprintf(":%s", cfg.APIServerPort)
		} else {
			log.Fatal().Msg("api_server_port option is deprecated; you cannot specify api_sever_port and api_server_bind at the same time.")
		}
	}
	if cfg.APIServerBind == "" {
		// defaults to :8080
		cfg.APIServerBind = ":8080"
	}

	return cfg
}

func main() {
	go initProfiler()

	cfg := parseCommandLineAndReadConfig()
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

	switch cfg.TLS {
	case false:
		if err := http.ListenAndServe(cfg.APIServerBind, r); err != nil {
			log.Fatal().Msgf("Failed to run HTTP server without TLS: %v", err)
		}
	case true:
		if err := http.ListenAndServeTLS(cfg.APIServerBind, cfg.ServerCertPath, cfg.ServerKeyPath, r); err != nil {
			log.Fatal().Msgf("Failed to run HTTP server with TLS: %v", err)
		}
	}
}
