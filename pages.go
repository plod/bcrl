package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

//TemplateData is a struct for passing through to templates
type TemplateData struct {
	PageTitle string
	Error     bool
	OtherData interface{}
	Year      int
}

var pushPromises = [...]string{"/assets/bcrl.css", "/assets/favicon.png"}

func notImplemented() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		generateError(w, http.StatusNotImplemented)
	}
}

func notFound(w http.ResponseWriter, req *http.Request) {
	generateError(w, http.StatusNotFound)
}

func home(w http.ResponseWriter, req *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	generateHTML(w, TemplateData{"", false, "The time is: " + tm, 1974}, "layout.tpl")
}

func generateError(w http.ResponseWriter, code int) {
	td := TemplateData{
		fmt.Sprintf("%v %v", code, http.StatusText(code)),
		true,
		"",
		0,
	}
	w.WriteHeader(code)
	generateHTML(w, td, "layout.tpl")
}

func generateHTML(w http.ResponseWriter, td TemplateData, filenames ...string) {
	if pusher, ok := w.(http.Pusher); ok {
		// Push is supported.
		for _, item := range pushPromises {
			if err := pusher.Push(item, nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
	}
	var files []string
	for _, file := range filenames {
		files = append(files, filepath.Join("templates", fmt.Sprintf("%s.html", file)))
	}
	templates := template.Must(template.ParseFiles(files...))
	td.Year = time.Now().Year()
	err := templates.ExecuteTemplate(w, "layout", td)
	if err != nil {
		log.Println(yellow("WARN"), err)
	}
}
