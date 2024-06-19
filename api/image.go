package api

import (
	"bytes"
	"errors"
	"image"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"www/internals"

	"github.com/chai2010/webp"
	"github.com/sunshineplan/imgconv"
)

func ImgOptimize(i image.Image, threshold int) image.Image {
	w := i.Bounds().Max.X

	if threshold >= w {
		return i
	}

	d := w / threshold
	return imgconv.Resize(i, &imgconv.ResizeOption{Width: w / d})
}

func Image(w http.ResponseWriter, r *http.Request) {
	error := internals.HttpErrorHelper(w)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if error("Error trying to parse query parameters", err, http.StatusInternalServerError) {
		return
	}

	if _, some := params["url"]; !some {
		error("\"url\" parameter missing", errors.New("Missing argument"), http.StatusBadRequest)
		return
	}
	u, err := url.Parse(params.Get("url"))
	if error("\"url\" is not a valid URL string", err, http.StatusBadRequest) {
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

	if _, some := params["threshold"]; !some {
		error("\"threshold\" parameter missing", errors.New("Missing argument"), http.StatusBadRequest)
		return
	}
	threshold, err := strconv.Atoi(params.Get("threshold"))
	if error("\"threshold\" parameter is not a valid integer", err, http.StatusBadRequest) {
		return
	}

	res, err := http.Get(u.String())
	if error("Error trying to fetch the image", err, http.StatusInternalServerError) {
		return
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		error(
			"Error trying to fetch the image, response is a non 2XX code",
			errors.New("Status code: "+res.Status),
			http.StatusInternalServerError,
		)
	}

	data, err := io.ReadAll(res.Body)
	if error("Error trying to read the image data", err, http.StatusInternalServerError) {
		return
	}

	img, err := imgconv.Decode(bytes.NewReader(data))
	if error("Error trying to decode the image", err, http.StatusInternalServerError) {
		return
	}

	img = ImgOptimize(img, threshold)

	err = webp.Encode(w, img, &webp.Options{Lossless: true})
	if error("Error trying to encode the image", err, http.StatusInternalServerError) {
		return
	}

	w.Header().Add("Cache-Control", "max-age=604800, stale-while-revalidate=86400, stale-if-error=86400")
	w.Header().Add("CDN-Cache-Control", "max-age=604800")
	w.Header().Add("Content-Type", "image/webp")
}
