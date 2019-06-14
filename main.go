package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	history = readHistory()

	templates["view"] = template.Must(template.New("").
		Funcs(template.FuncMap{"markdownToHTML": markdownToHTML}).
		ParseFiles("templates/view.html", "templates/base.html"))
	templates["edit"] = template.Must(
		template.ParseFiles("templates/edit.html", "templates/base.html"))

	staticFileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))
	mediaFileServer := http.FileServer(http.Dir("media"))
	http.Handle("/media/", http.StripPrefix("/media/", mediaFileServer))

	http.HandleFunc("/", redirectToHomeHandler)
	http.HandleFunc("/pages/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
