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
	"gopkg.in/natefinch/lumberjack.v2"
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
	logger *Logger
	// Logfile is the file name for this application log. The default is "ogp-app.log".
	Logfile = "/var/log/ogp-app.log"
	// MaxLogSize is the maximum log file size in MB before log lotation. The default is 10.
	MaxLogSize = 10
	// MaxLogBackups is the maximum number of lotated log files. The default is 5.
	MaxLogBackups = 5
	// MaxLogAge is the max age of each log file in days. The default is 28 days.
	MaxLogAge = 28
)

// Logger is the custom logger from rs/zerolog that works with natefinch/lumberjack for log lotation.
type Logger struct {
	*zerolog.Logger
}

func init() {
	// if in debug mode, the app emits the log to stdout as well as the log file.
	logger = initLogger(*debug)
}

// initLogger returns the common logger inside the application.
func initLogger(stdout bool) *Logger {
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

	w := &lumberjack.Logger{
		Filename:   Logfile,
		MaxSize:    MaxLogSize,
		MaxBackups: MaxLogBackups,
		MaxAge:     MaxLogAge,
	}
	var writers []io.Writer
	writers = append(writers, w)

	if stdout {
		writers = append(writers, os.Stdout)
	}

	mw := io.MultiWriter(writers...)

	l := zerolog.New(mw).With().Timestamp().Str("logName", logname).Logger()
	return &Logger{
		Logger: &l,
	}
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
