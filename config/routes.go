package config

import (
	"github.com/a-h/templ"
	"math/rand"
	"net/http"

	"www/api"
	"www/components"
	"www/internals"
	"www/pages"
)

var mockProjects = []components.Project{
	{
		Name:     "rec-sh",
		Summary:  "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link:     templ.SafeURL("https://github.com/dot013/rec-sh"),
		Image:    templ.SafeURL("https://images.unsplash.com/photo-1461749280684-dccba630e2f6?q=80&w=2669&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"),
		Icon:     "i-solar:programming-bold",
		WIP:      false,
		Current:  false,
		Language: "bash",
	},
	{
		Name:     ".mdparser",
		Summary:  "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link:     templ.SafeURL("https://github.com/dot013/rec-sh"),
		Image:    templ.SafeURL("https://images.unsplash.com/photo-1560697024-fd4affa63094?q=80&w=2630&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"),
		Icon:     "i-simple-icons:rust",
		WIP:      true,
		Current:  false,
		Language: "rust",
	},
	{
		Name:     ".www",
		Summary:  "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link:     templ.SafeURL("https://github.com/dot013/.www"),
		Image:    templ.SafeURL("https://images.unsplash.com/photo-1461749280684-dccba630e2f6?q=80&w=2669&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"),
		Icon:     "i-simple-icons:go",
		WIP:      true,
		Current:  true,
		Language: "golang",
	},
	{
		Name:     "Project",
		Summary:  "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link:     templ.SafeURL("https://github.com/dot013/.www"),
		Image:    templ.SafeURL("https://images.unsplash.com/photo-1461749280684-dccba630e2f6?q=80&w=2669&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"),
		Icon:     "i-solar:box-bold-duotone",
		WIP:      false,
		Current:  false,
		Language: "nix",
	},
	{
		Name:     "Project",
		Summary:  "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link:     templ.SafeURL("https://github.com/dot013/.www"),
		Image:    templ.SafeURL("https://images.unsplash.com/photo-1461749280684-dccba630e2f6?q=80&w=2669&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"),
		Icon:     "i-solar:box-bold-duotone",
		WIP:      false,
		Current:  false,
		Language: "nix",
	},
}

var mockBlog = []components.Blog{
	{
		Title: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim.",
		Summary: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint " +
			"cillum sint consectetur cupidatat Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link: templ.URL("/001"),
	},
	{
		Title: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim.",
		Summary: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint " +
			"cillum sint consectetur cupidatat Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link: templ.URL("/001"),
	},
	{
		Title: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim.",
		Summary: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint " +
			"cillum sint consectetur cupidatat Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link: templ.URL("/001"),
	},
	{
		Title: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim.",
		Summary: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint " +
			"cillum sint consectetur cupidatat Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link: templ.URL("/001"),
	},
	{
		Title: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim.",
		Summary: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint " +
			"cillum sint consectetur cupidatat Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link: templ.URL("/001"),
	},
	{
		Title: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim.",
		Summary: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint " +
			"cillum sint consectetur cupidatat Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		Link: templ.URL("/001"),
	},
}

var images = []string{
	"/images/image-1.png",
	"/images/image-2.png",
	"/images/image-3.png",
	"/images/image-4.png",
	"/images/image-5.png",
	"/images/image-6.png",
}

func init() {
	rand.Shuffle(len(images), func(i, j int) {
		images[i], images[j] = images[j], images[i]
	})
}

var ROUTES = []internals.Page{
	{Path: "index.html", Component: pages.Homepage(pages.HomepageProps{
		Projects: mockProjects,
		Images:   []string{},
		Blogs:    mockBlog,
	})},
}

func APIROUTES(mux *http.ServeMux) {
	mux.HandleFunc("/api/image", api.Image)
	mux.HandleFunc("/robots.txt", api.RobotsTxt)
	mux.HandleFunc("/ai.txt", api.AiTxt)
}
