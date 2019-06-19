package main

import (
	"log"
	"net/http"
)

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
	http.HandleFunc("/timers/", timersHandler)
	http.HandleFunc("/tools/", redirectToCardgenHandler)
	http.HandleFunc("/tools/cardgen/", cardgenHandler)
	http.HandleFunc("/tools/compare/", compareHandler)
	http.HandleFunc("/tools/deduplicate/", deduplicateHandler)
	http.HandleFunc("/tools/match/", matchHandler)
	http.HandleFunc("/uploads/", uploadHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
