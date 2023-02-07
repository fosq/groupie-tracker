package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

type IO struct {
	Input  string
	Output string
}

type Endpoints struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

type Artists []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Locations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

func main() {
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("templates"))))

	http.HandleFunc("/", formHandler) // Handles /ascii-art

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

	w.WriteHeader(http.StatusOK)
	whtml.Execute(w, io)

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(responseData))

	var result Artists
	json.Unmarshal(responseData, &result)

	io.Output = result[0].Name

	whtml, err = template.ParseFiles("templates/artistBubble.html")
	if err != nil {
		http.Error(w, "404 - Resource not found", http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
	whtml.Execute(w, io)
}

/*func api() {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(responseData))

	var result Artists
	json.Unmarshal(responseData, &result)

	//fmt.Println(reflect.TypeOf(&result))
	fmt.Println(result[0].Name)
	fmt.Println(len(result[0].Members))
	fmt.Println(len(result))

	// Print all artists
	for i := 0; i < len(result); i++ {
		fmt.Println(result[i].Name)
	}
	fmt.Println()

	fmt.Printf("Name: %v\nMembers: %v\nCreation Date: %v\nFirst album: %v\n", result[0].Name, result[0].Members, result[0].CreationDate, result[0].FirstAlbum)
}*/
