package main

import (
    "net/http"
    "regexp"
    "strings"
)

var validPath = regexp.MustCompile("^/(edit|save|pages)/([-a-zA-Z0-9]+)$")

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