package main

import (
	"image"
	"image/draw"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type RenderSpec struct {
	Width, Height int
	FontSize      float64
	Font          *truetype.Font
}

func (rs *RenderSpec) Render(text string) (image.Image, error) {
	logger.Info().Str("words", text).Send()
	img := image.NewRGBA(image.Rect(0, 0, rs.Width, rs.Height))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	opt := truetype.Options{
		Size: rs.FontSize,
	}
	face := truetype.NewFace(rs.Font, &opt)
	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	x := (fixed.I(rs.Width) - dr.MeasureString(text)) / 2
	dr.Dot.X = x
	y := (rs.Height + int(rs.FontSize)/2) / 2
	dr.Dot.Y = fixed.I(y)

	dr.DrawString(text)

	return img, nil
}
