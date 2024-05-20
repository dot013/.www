package main

import (
	"context"
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"guz.one/pages"
)

const PERMISSIONS = 0755

type Page struct {
	path      string
	component templ.Component
}

type Writer struct {
	root    *string
	pages   []Page
	context context.Context
}

func (w Writer) writeFile(path string, writer func(ctx context.Context, w io.Writer) error) {
	directory := filepath.Dir(path)
	err := os.MkdirAll(directory, PERMISSIONS)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = writer(w.context, f)
	if err != nil {
		log.Fatal(err)
	}
}

func (w Writer) WriteAll() {
	for _, page := range w.pages {
		p := filepath.Join(*w.root, page.path)
		log.Printf("Writing page %s", p)
		w.writeFile(p, page.component.Render)
	}
}

func main() {
	dir := flag.String("d", "./dist", "the directory to write the files")
	staticDir := flag.String("s", "./static", "the directory to copy static files from")

	pages := []Page{
		{"index.html", pages.Homepage()},
	}
	w := Writer{dir, pages, context.Background()}
	w.WriteAll()

	a, err := filepath.Abs(*staticDir)
	if err != nil {
		log.Fatal(err)
	}

	copyFile := func(root string, dest *string) func(string, fs.DirEntry, error) error {
		return func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			} else if path == root {
				return nil
			}

			c, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			p := strings.TrimPrefix(path, root)
			log.Printf("Copying static file %s to %s directory", p, *dest)

			p = filepath.Join(*dest, p)

			err = os.MkdirAll(filepath.Dir(p), PERMISSIONS)
			if err != nil {
				return err
			}

			f, err := os.Create(p)
			if err != nil {
				return err
			}
			defer f.Close()

			b, err := f.Write(c)
			if err != nil {
				return err
			}
			log.Printf("Wrote %v bytes in %s", b, p)

			return nil
		}
	}
	err = filepath.WalkDir(a, copyFile(a, dir))
	if err != nil {
		log.Fatal(err)
	}

}
