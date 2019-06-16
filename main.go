package main

import (
	"log"
	"net/http"
	"os"
)

// TemplateData context used to render templates
type TemplateData struct {
	Page          *Page
	History       *[]string
	Sitemap       *[]os.FileInfo
	SearchResults *SearchResults
}

var templateData TemplateData

const port = ":7713"

func main() {
	staticFileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))
	mediaFileServer := http.FileServer(http.Dir("media"))
	http.Handle("/media/", http.StripPrefix("/media/", mediaFileServer))

	http.HandleFunc("/", redirectToHomeHandler)
	http.HandleFunc("/pages/", viewPageHandler)
	http.HandleFunc("/edit/", editPageHandler)
	http.HandleFunc("/save/", savePageHandler)
	http.HandleFunc("/sitemap/", sitemapHandler)
	http.HandleFunc("/search/", searchHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
