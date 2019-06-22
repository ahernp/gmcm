package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var version = "0.6.0"

func main() {
	var port = flag.String("port", "7713", "Local port to listen on")
	flag.Parse()
	serve(*port)
}

func serve(port string) {
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
	http.HandleFunc("/tools/match/", matchHandler)
	http.HandleFunc("/uploads/", uploadHandler)

	fmt.Printf("Version %s; Listening on port :%s\n", version, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
