package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// UploadTemplateData template context
type UploadTemplateData struct {
	UploadedFiles *[]UploadedFile
	MainMenu      *string
	History       *[]string
}

// UploadedFile template data
type UploadedFile struct {
	Dir  string
	File os.FileInfo
}

const postMethod = "POST"

var uploadDirs = [...]string{"img", "code", "doc", "thumb"}
var uploadedFiles = getUploadedFiles()

var uploadTemplate = template.Must(
	template.ParseFiles("templates/uploads.html", "templates/base.html"))

func getUploadedFiles() []UploadedFile {
	var uploadedFiles []UploadedFile
	for dirPos := 0; dirPos < len(uploadDirs); dirPos++ {
		path := "media/" + uploadDirs[dirPos]
		files, _ := ioutil.ReadDir(path)
		for filePos := 0; filePos < len(files); filePos++ {
			uploadedFiles = append(uploadedFiles, UploadedFile{Dir: uploadDirs[dirPos], File: files[filePos]})
		}
	}
	return uploadedFiles
}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == postMethod {
		dir := request.FormValue("dir")
		request.ParseMultipartForm(10 << 20)
		sourceFile, handler, err := request.FormFile("newFile")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer sourceFile.Close()
		destFile, err := os.OpenFile("media/"+dir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer destFile.Close()
		io.Copy(destFile, sourceFile)
		uploadedFiles = getUploadedFiles()
	}
	templateData := UploadTemplateData{UploadedFiles: &uploadedFiles, MainMenu: &mainMenu, History: &history}
	uploadTemplate.ExecuteTemplate(writer, "base", templateData)
}
