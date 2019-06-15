package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

const port = ":7713"

var sitemap []os.FileInfo

func main() {
	history = readHistory()
	sitemap, _ = listPages()

	templates["view"] = template.Must(template.New("").
		Funcs(template.FuncMap{"markdownToHTML": markdownToHTML}).
		ParseFiles("templates/view.html", "templates/base.html"))
	templates["edit"] = template.Must(
		template.ParseFiles("templates/edit.html", "templates/base.html"))
	templates["sitemap"] = template.Must(
		template.ParseFiles("templates/sitemap.html", "templates/base.html"))
	templates["search"] = template.Must(
		template.ParseFiles("templates/search.html", "templates/base.html"))

	staticFileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))
	mediaFileServer := http.FileServer(http.Dir("media"))
	http.Handle("/media/", http.StripPrefix("/media/", mediaFileServer))

	http.HandleFunc("/", redirectToHomeHandler)
	http.HandleFunc("/pages/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/sitemap/", sitemapHandler)
	http.HandleFunc("/search/", searchHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}
