package main

import (
	"log"
	"net/http"
	"os"
)

// TemplateData context used to render all templates
type TemplateData struct {
	Page            *Page
	History         *[]string
	Sitemap         *[]os.FileInfo
	SearchResults   *SearchResults
	UploadedFiles   *[]UploadedFile
	CardgenData     *CardgenData
	CompareData     *CompareData
	DeduplicateData *DeduplicateData
}

var templateData TemplateData

const port = ":7713"

func main() {
	mediaFileServer := http.FileServer(http.Dir("media"))
	http.Handle("/media/", http.StripPrefix("/media/", mediaFileServer))
	staticFileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	http.HandleFunc("/", redirectToHomeHandler)
	http.HandleFunc("/edit/", editPageHandler)
	http.HandleFunc("/pages/", viewPageHandler)
	http.HandleFunc("/save/", savePageHandler)
	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/sitemap/", sitemapHandler)
	http.HandleFunc("/tools/", redirectToCardgenHandler)
	http.HandleFunc("/tools/cardgen/", cardgenHandler)
	http.HandleFunc("/tools/compare/", compareHandler)
	http.HandleFunc("/tools/deduplicate/", deduplicateHandler)
	http.HandleFunc("/uploads/", uploadHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
