package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/profiler"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	serviceName    = "ogp-app"
	serviceVersion = "1.0.0" // TODO: replace this
	maxRetry       = 3
)

var (
	// debug is the flag to change default log level to debug.
	debug = flag.Bool("debug", false, "Set log level to debug")

	// logger is the common logger among the apprecation.
	// TODO: confirm the performance with Cloud Profiler.
	logger zerolog.Logger
)

func init() {
	logger = initLogger(nil)
}

// initLogger returns the common logger inside the application.
func initLogger(w io.Writer) (logger zerolog.Logger) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// set up the field name and format to match Cloud Logging.
	// c.f. https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "severity"

	projectID, err := metadata.ProjectID()
	if err != nil {
		log.Warn().Msgf("Failed to retrieve Google Cloud Platform Project ID: %v", err)
		projectID = "local"
	}
	logname := fmt.Sprintf("projects/%s/logs/%s", projectID, serviceName)

	if w == nil {
		w = os.Stdout
	}
	return zerolog.New(w).With().Str("logName", logname).Logger()
}

// initProfiler starts Cloud Profiler. It retries at most maxRetry times.
func initProfiler() {
	for i := 0; i < maxRetry; i++ {
		// Profiler initialization, best done as early as possible.
		if err := profiler.Start(profiler.Config{
			Service:        serviceName,
			ServiceVersion: serviceVersion,
		}); err != nil {
			log.Printf("Failed to launch Profiler (%vth trial): %v-n", i, err)
		} else {
			log.Print("Started Cloud Profiler")
			return
		}
		d := time.Second * 5 * time.Duration(i)
		log.Printf("Wait %v seconds to retry launching Cloud Profiler\n", d)
		time.Sleep(d)
	}
	log.Printf("Failed to launch Cloud Profiler after %v trials\n", maxRetry)
}
