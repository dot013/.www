package api

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"mime"
	"net/http"
	"net/url"
	"strconv"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

func Scale(i image.Image, s float64) image.Image {
	r := i.Bounds()
	w, h := float64(r.Max.X), float64(r.Max.Y)

	var nw, nh int
	if s < 0 {
		s = s * -1
		nw, nh = int(math.Round(w/s)), int(math.Round(h/s))
	} else if s > 0 {
		s = s * -1
		nw, nh = int(math.Round(w*s)), int(math.Round(h*s))
	} else {
		nw, nh = int(w), int(h)
	}

	return imaging.Resize(i, nw, nh, imaging.CatmullRom)
}

func Image(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	var u *url.URL
	if _, some := params["url"]; !some {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("\"url\" parameter missing"))
		return
	} else {
		u, err = url.Parse(params.Get("url"))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		if u.Hostname() == "" {
			if r.URL.Scheme == "" {
				u.Scheme = "https"
			} else {
				u.Scheme = r.URL.Scheme
			}
			u.Host = r.Host
		}
	}

	var scale, width, height int
	if _, some := params["scale"]; !some {
		if _, some := params["width"]; !some {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("\"width\" parameter missing"))
			return
		} else {
			width, err = strconv.Atoi(params.Get("width"))

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("\"width\" parameter is not a valid integer"))
				return
			}
		}
		if _, some := params["height"]; !some {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("\"width\" parameter missing"))
			return
		} else {
			height, err = strconv.Atoi(params.Get("height"))

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("\"height\" parameter is not a valid integer"))
				return
			}
		}
	} else {
		scale, err = strconv.Atoi(params.Get("scale"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("\"scale\" parameter is not a valid integer"))
			return
		}
	}

	imgRes, err := http.Get(u.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if imgRes.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Url returned a status code of %v", imgRes.StatusCode)))
		return
	}

	data, err := io.ReadAll(imgRes.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	mime := mime.TypeByExtension(u.String())
	if mime == "" {
		mime = http.DetectContentType(data)
	}

	var img image.Image
	reader := bytes.NewReader(data)
	switch mime {
	case "image/png":
		img, err = png.Decode(reader)
	case "image/jpeg":
		img, err = jpeg.Decode(reader)
	case "image/gif":
		img, err = gif.Decode(reader)
	case "image/webp":
		img, err = webp.Decode(reader)
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("image is not either of \"jpeg\", \"png\", \"gif\", or \"webp\""))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error decoding the image:\n" + err.Error()))
		return
	}

	if scale != 0 {
		img = Scale(img, float64(scale))
	} else if width > 0 && height > 0 {
		img = imaging.Resize(img, width, height, imaging.CatmullRom)
	}

	switch mime {
	case "image/png":
		err = png.Encode(w, img)
	case "image/jpeg":
		err = jpeg.Encode(w, img, &jpeg.Options{Quality: 70})
	case "image/gif":
		err = gif.Encode(w, img, &gif.Options{NumColors: 256})
	case "image/webp":
		err = webp.Encode(w, img, &webp.Options{Lossless: true})
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("\"type\" parameter is not either of \"jpeg\", \"png\", \"gif\", or \"webp\""))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error encoding the image:\n" + err.Error()))
		return
	}

	w.Header().Add("Cache-Control", "max-age=604800, stale-while-revalidate=86400, stale-if-error=86400")
	w.Header().Add("CDN-Cache-Control", "max-age=604800")
	w.Header().Add("Content-Type", mime)
}
