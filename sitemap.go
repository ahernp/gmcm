package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var sitemap, _ = listPages()

var sitemapTemplate = template.Must(
	template.ParseFiles("templates/sitemap.html", "templates/base.html"))

func listPages() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(pagesPath)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return files, err
}

func renderSitemapTemplate(w http.ResponseWriter, sitemap *[]os.FileInfo) error {
	templateData = TemplateData{Sitemap: sitemap, History: &history}
	return sitemapTemplate.ExecuteTemplate(w, "base", templateData)
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	files, err := listPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sitemap = files
	renderSitemapTemplate(w, &sitemap)
}
