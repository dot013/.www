package internals

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

type Image struct {
	image.Image
	mime        string
	JpegOptions jpeg.Options
	GifOptions  gif.Options
	WebpOptions webp.Options
}

func NewImage(b []byte) (Image, error) {
	m := http.DetectContentType(b)

	r := bytes.NewReader(b)

	var img image.Image
	var err error

	switch m {
	case "image/png":
		img, err = png.Decode(r)
	case "image/jpeg":
		img, err = jpeg.Decode(r)
	case "image/gif":
		img, err = gif.Decode(r)
	case "image/webp":
		img, err = webp.Decode(r)
	default:
		err = errors.ErrUnsupported
	}
	if err != nil {
		return Image{}, err
	}

	image := img.(Image)
	image.mime = m
	image.JpegOptions = jpeg.Options{
		Quality: 70,
	}
	image.GifOptions = gif.Options{
		NumColors: 256 / 2,
	}
	image.WebpOptions = webp.Options{
		Lossless: true,
		Quality:  70.0,
		Exact:    true,
	}

	return image, nil

}

func (i *Image) Decode(b []byte) error {
	img, err := NewImage(b)
	*i = img
	return err
}

func (i *Image) Encode(w io.Writer) error {
	var err error

	switch i.mime {
	case "image/png":
		err = png.Encode(w, i.Image)
	case "image/jpeg":
		err = jpeg.Encode(w, i.Image, &jpeg.Options{Quality: 70})
	case "image/gif":
		err = gif.Encode(w, i.Image, &gif.Options{NumColors: 256})
	case "image/webp":
		err = webp.Encode(w, i.Image, &webp.Options{Lossless: true})
	default:
		err = errors.ErrUnsupported
	}
	if err != nil {
		return err
	}

	return nil
}

func (i *Image) Quality(q float64) {
	i.WebpOptions.Quality = float32(q)
	i.JpegOptions.Quality = int(math.Round(q))
}

func (i *Image) Optimize(threshold int) {
	w := i.Bounds().Max.X

	if threshold >= w {
		return
	}

	d := threshold / w
	i.Scale(d * -1)
}

func (i *Image) Scale(s int) {
	r := i.Bounds()
	w, h := r.Max.X, r.Max.Y

	var nw, nh int
	if s < 0 {
		s = s * -1
		nw, nh = w/s, h/s
	} else if s > 0 {
		s = s * -1
		nw, nh = w*s, h*s
	} else {
		nw, nh = w, h
	}

	imaging.Resize(i, nw, nh, imaging.CatmullRom)
}

func (i *Image) GetMime() string {
	return i.mime
}
