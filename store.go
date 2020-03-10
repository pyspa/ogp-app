package main

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"path/filepath"

	"cloud.google.com/go/storage"
)

const gcsHost = "storage.googleapis.com"

type ImageStore interface {
	Save(ctx context.Context, img image.Image, id string) error
	URL(id string) string
}

type LocalImageStore struct {
	BasePath      string
	PublicBaseURL string
}

func (s *LocalImageStore) Save(_ context.Context, img image.Image, id string) error {
	filename := fmt.Sprintf("%s.png", id)
	filepath := filepath.Join(s.BasePath, filename)

	outfile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	if err := png.Encode(outfile, img); err != nil {
		return err
	}
	return nil
}

func (s *LocalImageStore) URL(id string) string {
	filename := fmt.Sprintf("%s.png", id)
	return path.Join(s.PublicBaseURL, filename)
}

type GCSImageStore struct {
	Client        *storage.Client
	Bucket        string
	Prefix        string
	PublicBaseURL string
}

func (s *GCSImageStore) Save(ctx context.Context, img image.Image, id string) error {
	path := path.Join(s.Prefix, id)
	w := s.Client.Bucket(s.Bucket).Object(path).NewWriter(ctx)
	defer w.Close()
	w.ContentType = "image/png"
	if err := png.Encode(w, img); err != nil {
		return err
	}
	return nil
}

func (s *GCSImageStore) URL(id string) string {
	path := path.Join(s.Prefix, id)
	if s.PublicBaseURL == "" {
		return fmt.Sprintf("https://%s/%s/%s", gcsHost, s.Bucket, path)
	} else {
		return fmt.Sprintf("%s/%s", s.PublicBaseURL, path)
	}
}
