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

type Card struct {
	Id           int
	FirstAlbum   string
	Location     []string
	Image        string
	GroupName    string
	CreationDate int
	Members      []string
}

type Cards []Card

type MainData struct {
	Cards        []Card
	CountMembers []int
	GroupNames   []string
	CreationDate []int
	FirstAlbum   []string
	Members      []string
	Locations    []string
}

type Artist struct {
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

	whtml, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, "404 - Resource not found", http.StatusNotFound)
	}

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

	var result []Artist
	var cards Cards
	json.Unmarshal(responseData, &result)

	// Loop through each artist
	for _, artist := range result {
		data := Card{
			Id:           artist.ID,
			Image:        artist.Image,
			GroupName:    artist.Name,
			CreationDate: artist.CreationDate,
			FirstAlbum:   artist.FirstAlbum,
			Members:      artist.Members,
		}

		cards = append(cards, data)
	}

	//whtml, err = template.ParseFiles("templates/artistBubble.html")
	//if err != nil {
	//	http.Error(w, "404 - Resource not found", http.StatusNotFound)
	//}

	//w.WriteHeader(http.StatusOK)
	whtml.Execute(w, cards)
}

//    <a href="artist/{{.ID}}">
