package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var sitemap []os.FileInfo

func listPages() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(pagesPath)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return files, err
}

func renderSitemapTemplate(w http.ResponseWriter, name string, sitemap *[]os.FileInfo) error {
	template, ok := templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	templateData = TemplateData{Sitemap: sitemap, History: &history}
	return template.ExecuteTemplate(w, "base", templateData)
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	files, err := listPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sitemap = files
	renderSitemapTemplate(w, "sitemap", &sitemap)
}
