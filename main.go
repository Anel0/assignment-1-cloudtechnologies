package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type BookCountResponse struct {
	Language	string `json:"language"`
	Books	int	`json:"books"`
	Authors int	`json:"authors"`
	Fraction float32 `json:"fraction"`
}

type Person struct {
	BirthYear int `json:"birth_year"`
	DeathYear int `json:"death_year"`
	Name string `json:"name"`
}

type Book struct {
	Id    int    `json:"id"`
	Title	string	`json:"title"`
	Subjects	[]string	`json:"subjects"`
	Authors	[]Person	`json:"authors"`
	Translators	[]Person	`json:"translators"`
	Bookshelves	[]string	`json:"bookshelves"`
	Languages	[]string	`json:"languages"`
	Copyright	bool	`json:"copyright"`
	MediaType	map[string][]string	`json:"media_type"`
	Formats	string	`json:"formats"`
	DownloadCount	int	`json:"download_count"`
}

type GutenbergBookResponse struct {
	Count    int    `json:"count"`
	Next    string    `json:"next"`
	Previous    string    `json:"previous"`
	Books []Book `json:"results"`
}

func findUniqueAuthors(books []Book) int {
	authorMap := make(map[string]int)
	for i := 0; i < len(books); i++ {
		var authors = books[i].Authors
		for j := 0; j < len(authors); j++ {
			authorMap[authors[j].Name] = 1
		}
	}
	return len(authorMap)
}

func getTotalBooks() int {
	responseBooks, err := http.Get("http://129.241.150.113:8000/books")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(responseBooks.Body)
	if err != nil {
		log.Fatal(err)
	}
	var response GutenbergBookResponse
	json.Unmarshal(responseData, &response)
	return response.Count
}

func bookCount(w http.ResponseWriter, r *http.Request){

	var queryParameters = r.URL.Query()
	if !queryParameters.Has("language"){
		return
	}
	responseBooks, err := http.Get("http://129.241.150.113:8000/books/?languages=" + queryParameters.Get("language"))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(responseBooks.Body)
	if err != nil {
		log.Fatal(err)
	}
	var response GutenbergBookResponse
	json.Unmarshal(responseData, &response)

	var ourResponse = BookCountResponse{
		queryParameters.Get("language"),
		response.Count,
		findUniqueAuthors(response.Books),
		float32(response.Count) / float32(getTotalBooks()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ourResponse)

	fmt.Println("Language:", queryParameters.Get("language"))
	fmt.Println("Books:", response.Count)
	fmt.Println("Authors:",findUniqueAuthors(response.Books))
	fmt.Println("Fraction:",float32(response.Count) / float32(getTotalBooks()))
	//fmt.Fprintf(w, string(responseData))
}

func readership(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Readership API")
	fmt.Println("Endpoint Hit: Readership")
}


func status(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Status API")
	fmt.Println("Endpoint Hit: Status")
}


func handleRequests() {
	http.HandleFunc("/librarystats/v1/bookcount/", bookCount)
	http.HandleFunc("/librarystats/v1/readership/", readership)
	http.HandleFunc("/librarystats/v1/status/", status)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func main() {
	handleRequests()
}
