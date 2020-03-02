package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/achiku/mux"
	"cloud.google.com/go/profiler"
)

const (
	serviceName = "ogp-app"
	serviceVersion = "1.0.0" // TODO: replace this
	maxRetry = 3
)

func init() {

}

func main() {
	go initProfiler()

	configFile := flag.String("c", "", "configuration file path")
	flag.Parse()

	cfg, err := NewConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to create a new configuration: %v", err)
	}
	app, err := NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create a new application: %v", err)
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
			log.Fatal(err)
		}
	case true:
		if err := http.ListenAndServeTLS(p, cfg.ServerCertPath, cfg.ServerKeyPath, r); err != nil {
			log.Fatal(err)
		}
	}
}


func initProfiler() {
	for i := 0; i < maxRetry; i++ {
		log.Println("")
		// Profiler initialization, best done as early as possible.
		if err := profiler.Start(profiler.Config{
			Service:        serviceName,
			ServiceVersion: serviceVersion,
		}); err != nil {
			log.Println("Failed to launch Profiler (%vth trial): %v", i, err)
		} else {
			log.Println("Started Cloud Profiler")
			return
		}
		d := time.Second * 5 * time.Duration(i)
		log.Printf("Wait %v seconds to retry launching Cloud Profiler\n", d)
		time.Sleep(d)
	}
	log.Printf("Failed to launch Cloud Profiler after %v trials\n", maxRetry)
}