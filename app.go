package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"cloud.google.com/go/storage"
	"github.com/BurntSushi/toml"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

// Config ogp.app config
type Config struct {
	BaseURL            string  `toml:"base_url"`
	APIServerBind      string  `toml:"api_server_bind"`
	APIServerPort      string  `toml:"api_server_port"`
	TLS                bool    `toml:"tls"`
	KoruriBoldFontPath string  `toml:"koruri_bold_font_path"`
	DefaultImageWidth  int     `toml:"default_image_width"`
	DefaultImageHeight int     `toml:"default_image_height"`
	DefaultFontSize    float64 `toml:"default_font_size"`
	ServerCertPath     string  `toml:"server_cert_path"`
	ServerKeyPath      string  `toml:"server_key_path"`
	StorageBackend     string  `toml:"storage_backend"`
	Storage            struct {
		Local struct {
			BasePath      string `toml:"base_path"`
			PublicBaseURL string `toml:"public_base_url"`
		} `toml:"local"`
		GCS struct {
			CredentialsFile string `toml:"credentials_file"`
			Bucket          string `toml:"bucket"`
			Prefix          string `toml:"prefix"`
			PublicBaseURL   string `toml:"public_base_url"`
		} `toml:"gcs"`
	} `toml:"storage"`
}

// App ogp.app
type App struct {
	Config        *Config
	KoruriBold    *truetype.Font
	store         ImageStore
	OgpPagePath   string
	IndexPagePath string
	OgpPageTmpl   *template.Template
	IndexPageTmpl string
}

// NewConfig create app config
func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open path %v: %w", path, err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %v: %w", path, err)
	}
	var cfg Config
	if err := toml.Unmarshal(buf, &cfg); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal toml data: %w", err)
	}
	return &cfg, nil
}

var storageBackends = map[string]func(context.Context, *Config) (ImageStore, error){
	"local": func(_ context.Context, cfg *Config) (ImageStore, error) {
		bcfg := cfg.Storage.Local
		if bcfg.BasePath == "" && bcfg.PublicBaseURL == "" {
			return &LocalImageStore{
				BasePath:      "data",
				PublicBaseURL: "/image",
			}, nil
		} else {
			return &LocalImageStore{
				BasePath:      bcfg.BasePath,
				PublicBaseURL: bcfg.PublicBaseURL,
			}, nil
		}
	},
	"gcs": func(ctx context.Context, cfg *Config) (ImageStore, error) {
		bcfg := cfg.Storage.GCS
		if bcfg.Bucket == "" {
			fmt.Printf("%+v", cfg)
			return nil, fmt.Errorf("configuration for GCS storage backend is missing")
		}
		var options []option.ClientOption
		if bcfg.CredentialsFile != "" {
			options = append(options, option.WithCredentialsFile(bcfg.CredentialsFile))
		}
		client, err := storage.NewClient(ctx, options...)
		if err != nil {
			return nil, fmt.Errorf("failed to create GCS client instance: %w", err)
		}
		return &GCSImageStore{
			Client:        client,
			Bucket:        bcfg.Bucket,
			Prefix:        bcfg.Prefix,
			PublicBaseURL: bcfg.PublicBaseURL,
		}, nil
	},
}

func buildStorageBackend(ctx context.Context, cfg *Config) (ImageStore, error) {
	backendName := cfg.StorageBackend
	if backendName == "" {
		backendName = "local"
	}
	factory, ok := storageBackends[backendName]
	if !ok {
		return nil, fmt.Errorf("unsupported storage backend: %s", cfg.StorageBackend)
	}
	return factory(ctx, cfg)
}

// NewApp create app
func NewApp(ctx context.Context, cfg *Config) (*App, error) {
	fontBytes, err := ioutil.ReadFile(cfg.KoruriBoldFontPath)
	if err != nil {
		return nil, err
	}
	ft, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	var ogpPageTmpl *template.Template
	{
		path := filepath.Join("client", "dist", "p.html")
		pbuf, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", path, err)
		}
		ogpPageTmpl, err = template.New("page").Parse(string(pbuf))
		if err != nil {
			return nil, err
		}
	}

	var indexPageTmpl string
	{
		path := filepath.Join("client", "dist", "index.html")
		idxbuf, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", path, err)
		}
		indexPageTmpl = string(idxbuf)
	}

	store, err := buildStorageBackend(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage backend %s: %w", cfg.StorageBackend, err)
	}

	return &App{
		Config:        cfg,
		KoruriBold:    ft,
		OgpPageTmpl:   ogpPageTmpl,
		IndexPageTmpl: indexPageTmpl,
		store:         store,
	}, nil
}

// deduceBaseURL deduces the application's base URL from the request headers.
func deduceBaseURL(r *http.Request) string {
	var scheme string
	proto := r.Header.Get("X-Forwarded-Proto")
	if proto != "" {
		switch proto {
		case "http", "https":
			scheme = proto + ":"
		}
	} else {
		if r.TLS != nil {
			scheme = "https:"
		} else {
			scheme = "http:"
		}
	}
	host := r.Host
	return fmt.Sprintf("%s//%s", scheme, host)
}

// BaseURL returns the base URL for the service
func (app *App) BaseURL(r *http.Request) string {
	if app.Config.BaseURL != "" {
		return app.Config.BaseURL
	} else {
		return deduceBaseURL(r)
	}
}

// OgpPage display ogp page
func (app *App) OgpPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	url := app.store.URL(id)
	data := map[string]string{
		"id":       id,
		"pageURL":  fmt.Sprintf("%s/p/%s", app.BaseURL(r), id),
		"imageURL": url,
	}
	w.WriteHeader(http.StatusOK)
	if err := app.OgpPageTmpl.Execute(w, data); err != nil {
		return
	}
	return
}

// IndexPage display index page
func (app *App) IndexPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, app.IndexPageTmpl)
	return
}

type createImageReq struct {
	Words string `json:"words"`
}

// CreateImage create ogp image API
func (app *App) CreateImage(w http.ResponseWriter, r *http.Request) {
	var d createImageReq
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		logger.Error().Msgf("decode failed: %s", err)
		return
	}
	words := d.Words
	id := uuid.New().String()

	img, err := (&RenderSpec{
		Width:    app.Config.DefaultImageWidth,
		Height:   app.Config.DefaultImageHeight,
		FontSize: app.Config.DefaultFontSize,
		Font:     app.KoruriBold,
	}).Render(words)
	if err != nil {
		logger.Error().Msgf("failed to render image: %w", err)
		return
	}

	err = app.store.Save(r.Context(), img, id)
	if err != nil {
		logger.Error().Msgf("create image failed: %s", err)
		return
	}

	url := app.store.URL(id)
	data := map[string]string{
		"words":    words,
		"id":       id,
		"pageURL":  fmt.Sprintf("%sp/%s", app.BaseURL(r), id),
		"imageURL": url,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		logger.Printf("encode failed: %s", err)
		return
	}
	return
}

// Redirector redirects to the image URL
func (app *App) Redirector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idAndMaybeExt := vars["id"]
	i := strings.LastIndexByte(idAndMaybeExt, '.')
	var id string
	if i >= 0 {
		id = idAndMaybeExt[:i]
	} else {
		id = idAndMaybeExt
	}

	url := app.store.URL(id)
	http.Redirect(w, r, url, http.StatusFound)
}

func (app *App) SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/").HandlerFunc(app.IndexPage)
	r.Methods(http.MethodGet).Path("/p/{id}").HandlerFunc(app.OgpPage)

	if s, ok := app.store.(*LocalImageStore); ok {
		r.Methods(http.MethodGet).PathPrefix("/image/").Handler(
			http.StripPrefix("/image/", http.FileServer(http.Dir(s.BasePath))))
	} else {
		r.Methods(http.MethodGet).Path("/image/{id}").Handler(http.HandlerFunc(app.Redirector))
	}

	// static asset
	r.Methods(http.MethodGet).PathPrefix("/js/").Handler(
		http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.Join("client", "dist", "js")))))
	r.Methods(http.MethodGet).PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join("client", "dist", "css")))))
	r.Methods(http.MethodGet).PathPrefix("/img/").Handler(
		http.StripPrefix("/img/", http.FileServer(http.Dir(filepath.Join("client", "dist", "img")))))

	// API
	r.Methods(http.MethodPost).Path("/api/image").Handler(
		loggingMiddleware(http.HandlerFunc(app.CreateImage)))

	return r
}
