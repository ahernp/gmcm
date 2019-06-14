package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Context struct {
	Page    *Page
	History []string
}

var context Context

const historySize = 20

var history = make([]string, historySize)

var templates = make(map[string]*template.Template)

var validPath = regexp.MustCompile("^/(edit|save|pages)/([-a-zA-Z0-9]+)$")

func readHistory() []string {
	filename := "data/history.txt"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return strings.Split(string(content), "\n")
}

func writeHistory() error {
	filename := "data/history.txt"
	historyAsString := strings.Join(history, "\n")
	return ioutil.WriteFile(filename, []byte(historyAsString), 0600)
}

func updateHistory(slug string) {
	newHistory := make([]string, 1, historySize)
	newHistory[0] = slug
	for i := 0; i < len(history); i++ {
		if history[i] != slug {
			newHistory = append(newHistory, history[i])
		}
	}
	history = newHistory[:historySize]
	writeHistory()
}

func renderTemplate(w http.ResponseWriter, name string, p *Page) error {
	template, ok := templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	context = Context{Page: p, History: history}
	return template.ExecuteTemplate(w, "base", context)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, slug string) {
	p, err := loadPage(slug)
	if err != nil {
		http.Redirect(w, r, "/edit/"+slug, http.StatusFound)
		return
	}
	updateHistory(slug)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, slug string) {
	p, err := loadPage(slug)
	if err != nil {
		p = &Page{Slug: slug}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, slug string) {
	body := r.FormValue("body")
	bodySansCarriageReturns := strings.ReplaceAll(body, "\r", "")
	p := &Page{Slug: slug, Body: []byte(bodySansCarriageReturns)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/pages/"+slug, http.StatusFound)
}

func redirectToHomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/pages/home", http.StatusFound)
}

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
