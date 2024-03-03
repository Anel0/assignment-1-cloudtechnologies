package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type BookCountResponse struct {
	Language	string `json:"language"`
	Books	int	`json:"books"`
	Authors int	`json:"authors"`
	Fraction float32 `json:"fraction"`
}

type ReadershipResponse struct {
	Country 	string	`json:"country"`
	Iso 		string	`json:"isocode"`
	Nooks		int		`json:"books"`
	Authors		int		`json:"authors"`
	Readerhship	int		`json:"readership"`

}


type Country struct {
	Iso_three 	string `json:"ISO3166_1_Alpha_3"`
	Iso_two			string `json:"ISO3166_1_Alpha_2"`
	OfficialName	string `json:"Official_Name"`
	Region		string	`json:"Region_Name"`
	SubRegion	string  `json:"Sub_Region_Name"`
	Language 	string	`json:"Language"`
}

type LanguageToCountriesResponse []Country

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

	var responseList []BookCountResponse
	var totalBooks = getTotalBooks()
	var languages = queryParameters.Get("language")
	var languageArray = strings.Split(languages, ",")


	for _, language := range languageArray{
		responseList = append(responseList, bookCountForSingleLanguage(language, totalBooks))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseList)


}

func bookCountForSingleLanguage(language string, totalBooks int) BookCountResponse{
	var page = 1
	var hasNextPage = true
	var allBooks []Book
	var count = 0

	for(hasNextPage){
		var url = fmt.Sprintf("http://129.241.150.113:8000/books/?languages=%s&page=%d",language,page)
		responseBooks, err:= http.Get(url)
		if err != nil{
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(responseBooks.Body)
		if err != nil {
			log.Fatal(err)
		}

		var response GutenbergBookResponse

		json.Unmarshal(responseData, &response)

		allBooks = append(allBooks, response.Books...)
		if response.Next == "" {
			hasNextPage = false
			count = response.Count
		}  else {
			page++
		}
	}


		var ourResponse = BookCountResponse{
			language,
			count,
			findUniqueAuthors(allBooks),
			float32(count) / float32(totalBooks),
		}

		return ourResponse
}

func readership(w http.ResponseWriter, r *http.Request){
	var language = path.Base(r.URL.Path)
	var countryNames = getCountriesFromLanguage(language,10)
	//

	for _, countryName := range countryNames{

	}

	fmt.Fprintf(w, "Readership API")
	fmt.Println("Endpoint Hit: Readership")

	//{
	//     "country": "Norway",
	//     "isocode": "NO",
	//     "books": 21,
	//     "authors": 14,
	//     "readership": 5379475
	//  },
}

func getCountriesFromLanguage(language string, limit int) []string{
	var url = fmt.Sprintf("http://129.241.150.113:3000/language2countries/%s",language)
	var counter = 0
	var countries []string
	responseCountries, err:= http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(responseCountries.Body)
	if err != nil {
		log.Fatal(err)
	}


	var response LanguageToCountriesResponse
	json.Unmarshal(responseData, &response)


	for _, country := range response{
		countries = append(countries, country.OfficialName)
		counter++
		if (counter > limit){
			return countries
		}
	}

	return countries

}


func getPopulationForCountry(countryName string){
	val url = fmt.Sprintf("https://restcountries.com/v3.1/name/%s",countryName)
	responseCountryData, err:= http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(responseCountryData.Body)
	if err != nil {
		log.Fatal(err)
	}

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
