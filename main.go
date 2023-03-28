package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type Artist struct {
	Id              int      `json:"id"`
	Image           string   `json:"image"`
	GroupName       string   `json:"name"`
	Members         []string `json:"members"`
	CreationDate    int      `json:"creationDate"`
	FirstAlbum      string   `json:"firstAlbum"`
	LocationsURL    string   `json:"locations"`
	ConcertDatesURL string   `json:"concertDates"`
	RelationsURL    string   `json:"relations"`
	Relations       map[string][]string
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var (
	link   string = "https://groupietrackers.herokuapp.com/api/artists"
	result []Artist
)

func main() {
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("templates"))))

	http.HandleFunc("/", formHandler) // Handles /ascii-art
	http.HandleFunc("/artist/", showArtistPage)

	fmt.Printf("Starting server at port 8080, access the page with 'localhost:8080' in a browser\n")
	fmt.Printf("Press 'Ctrl + C' to end the server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil { // Listens on port 8080
		log.Fatal(err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 Page Not Found.", http.StatusNotFound)
		return
	}
	whtml, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, "404 - Resource not found", http.StatusNotFound)
	}

	fetchData(link, 0)

	whtml.Execute(w, &result)
}

// Shows artist page when clicking on artist bubble
func showArtistPage(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`(^/artist/\d+)$`)
	if !re.MatchString(r.URL.Path) {
		http.Error(w, "404 Page not found.", http.StatusNotFound)
	}

	pageId := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(pageId)
	fetchData(link, id)

	whtml, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "404 - Resource not found", http.StatusNotFound)
	}
	whtml.Execute(w, result[id-1])

}

func fetchData(link string, id int) {
	response, err := http.Get(link)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseData, &result)
	if err != nil {
		fmt.Println(err)
	}

	if id != 0 {
		artist := &result[id-1]
		var relations Relations

		response, err := http.Get(artist.RelationsURL)
		if err != nil {
			log.Fatal(err)
		}
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(responseData, &relations)
		if err != nil {
			fmt.Println(err)
		}

		dateLocMap := make(map[string][]string)
		for key, val := range relations.DatesLocations {
			newKey := key
			newKey = strings.Replace(newKey, "-", ", ", -1)
			newKey = strings.Replace(newKey, "_", " ", -1)
			newKey = strings.Title(strings.ToLower(newKey))
			newKey = strings.Replace(newKey, "Usa", "USA", -1)
			newKey = strings.Replace(newKey, "Uk", "UK", -1)
			dateLocMap[newKey] = val
		}
		artist.Relations = dateLocMap
	}
}
