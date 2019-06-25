package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var version = "0.9.0"

// GlobalContext template data
type GlobalContext struct {
	MainMenu *string
	History  *[]string
	Version  *string
}

var globalContext = GlobalContext{MainMenu: &mainMenu, History: &history, Version: &version}

func main() {
	go cacheAllPages()
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
	http.HandleFunc("/uploads/", uploadHandler)

	fmt.Printf("Version %s; Listening on port :%s\n", version, port)
	fmt.Printf("Open http://localhost:%s to use\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
