package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type IO struct {
	Input  string
	Output string
}

type Artists struct {
}

func main() {

	http.Handle("/templates/", http.FileServer(http.Dir("templates")))

	http.HandleFunc("/", formHandler)

	fmt.Printf("Starting server at port 8080, access the page with 'localhost:8080' in a browser\n")
	fmt.Printf("Press 'Ctrl + C' to end the server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil { // Listens on port 8080
		log.Fatal(err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	var io IO
	whtml, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, "404 - Resource not found", http.StatusNotFound)
	}

	if r.URL.Path != "/" {
		http.Error(w, "404 - Page not found", http.StatusNotFound)
		return
	}

	if r.URL.Path == "/" {
		w.WriteHeader(http.StatusOK)
		whtml.Execute(w, io)
	}
}
