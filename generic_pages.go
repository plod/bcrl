package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func notImplemented() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func home(w http.ResponseWriter, req *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	generateHTML(w, "The time is: "+tm, "layout.tpl", "home.pge")
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, filepath.Join("templates", fmt.Sprintf("%s.html", file)))
	}
	templates := template.Must(template.ParseFiles(files...))
	err := templates.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		log.Println(yellow("WARN"), err)
	}
}
