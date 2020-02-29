package main

import (
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/achiku/mux"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Config ogp.app config
type Config struct {
	BaseURL            string  `toml:"base_url"`
	APIServerPort      string  `toml:"api_server_port"`
	KoruriBoldFontPath string  `toml:"koruri_bold_font_path"`
	DefaultImageWidth  int     `toml:"default_image_width"`
	DefaultImageHeight int     `toml:"default_image_height"`
	DefaultFontSize    float64 `toml:"default_font_size"`
	ServerCertPath     string  `toml:"server_cert_path"`
	ServerKeyPath      string  `toml:"server_key_path"`
}

// App ogp.app
type App struct {
	Config     *Config
	OutputTmpl *template.Template
	InputTmpl  string
	PageTmpl   *template.Template
	KoruriBold *truetype.Font
}

// NewConfig create app config
func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := toml.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// NewApp create app
func NewApp(cfg *Config) (*App, error) {
	outputTmpl, err := template.New("output").Parse(output)
	if err != nil {
		return nil, err
	}
	pageTmpl, err := template.New("page").Parse(page)
	if err != nil {
		return nil, err
	}
	fontBytes, err := ioutil.ReadFile(cfg.KoruriBoldFontPath)
	if err != nil {
		return nil, err
	}
	ft, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:     cfg,
		OutputTmpl: outputTmpl,
		InputTmpl:  inputForm,
		PageTmpl:   pageTmpl,
		KoruriBold: ft,
	}, nil
}

func isTLS(url string) bool {
	return strings.HasPrefix(url, "https")
}

func createImage(width, height int, fontsize float64, ft *truetype.Font, text, out string) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	opt := truetype.Options{
		Size: fontsize,
	}
	face := truetype.NewFace(ft, &opt)
	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	x := (fixed.I(width) - dr.MeasureString(text)) / 2
	dr.Dot.X = x
	y := (height + int(fontsize)/2) / 2
	dr.Dot.Y = fixed.I(y)

	dr.DrawString(text)

	outfile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outfile.Close()

	if err := png.Encode(outfile, img); err != nil {
		return err
	}
	return nil
}

const inputForm = `
<html>
  <head>
    <title>ogp.app</title>
  </head>
  <body>
    <h2><a href="/">ogp.app</a></h2>
	<form action="/image" method="post">
	  <p>
	    <textarea name="words" rows="1" size="20"></textarea>
	  </p>
	  <p>
	    <input type="submit" value="ogp">
	  </p>
	</form>
  </body>
</html>
`

const output = `
<html>
  <head>
    <title>ogp.app</title>
  </head>
  <body>
    <h2><a href="/">ogp.app</a></h2>
	<p>
	  {{ .words }}
	</p>
	<p>
	  <input size="60" value="{{ .baseURL }}/p/{{ .id }}">
	</p>
	<p>
	  <img src="/image/{{ .file }}"></img>
	</p>
  </body>
</html>
`

const page = `
<html>
  <head>
    <title>ogp.app</title>
	  <meta name="og:type" content="website" />
	  <meta name="og:url" content="{{ .baseURL }}/p/{{ .id }}" />
	  <meta name="og:title" content="ogp.app" />
	  <meta name="og:image" content="{{ .baseURL }}/image/{{ .file }}" />
	  <meta name="twitter:card" content="summary_large_image" />
	  <meta name="twitter:url" content="{{ .baseURL }}/p/{{ .id }}" />
  </head>
  <body>
    <h2><a href="/">ogp.app</a></h2>
	<p>
	  <img src="/image/{{ .file }}"></img>
	</p>
  </body>
</html>
`

// InputPage display input form
func (app *App) InputPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, inputForm)
	return
}

// OgpPage display ogp page
func (app *App) OgpPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	data := map[string]string{
		"id":      id,
		"file":    fmt.Sprintf("%s.png", id),
		"baseURL": app.Config.BaseURL,
	}
	w.WriteHeader(http.StatusOK)
	if err := app.PageTmpl.Execute(w, data); err != nil {
		return
	}
	return
}

// CreateImage create ogp image
func (app *App) CreateImage(w http.ResponseWriter, r *http.Request) {
	words := r.PostFormValue("words")
	id := uuid.New()
	filename := fmt.Sprintf("%s.png", id.String())
	filepath := path.Join("data", filename)
	wi, he, fs := app.Config.DefaultImageWidth, app.Config.DefaultImageHeight, app.Config.DefaultFontSize
	if err := createImage(wi, he, fs, app.KoruriBold, words, filepath); err != nil {
		return
	}
	log.Printf("post data: %s", words)
	w.WriteHeader(http.StatusOK)
	data := map[string]string{
		"words":   words,
		"file":    filename,
		"id":      id.String(),
		"baseURL": app.Config.BaseURL,
	}
	if err := app.OutputTmpl.Execute(w, data); err != nil {
		return
	}
	return
}
