package main

import (
	"log"
	"net/http"
	"os"
)

// TemplateData used when rendering templates
type TemplateData struct {
	Page          *Page
	History       *[]string
	Sitemap       *[]os.FileInfo
	SearchResults *SearchResults
}

var templateData TemplateData

const port = ":7713"
const pagesPath = "data/pages/"

func main() {
	history = readHistory()
	sitemap, _ = listPages()

	staticFileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))
	mediaFileServer := http.FileServer(http.Dir("media"))
	http.Handle("/media/", http.StripPrefix("/media/", mediaFileServer))

	http.HandleFunc("/", redirectToHomeHandler)
	http.HandleFunc("/pages/", makePageHandler(viewPageHandler))
	http.HandleFunc("/edit/", makePageHandler(editPageHandler))
	http.HandleFunc("/save/", makePageHandler(savePageHandler))
	http.HandleFunc("/sitemap/", sitemapHandler)
	http.HandleFunc("/search/", searchHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
