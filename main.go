package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

//structs-------------------------------------------------------
type BookCountResponse struct {
	Language	string  `json:"language"`
	Books		int		`json:"books"`
	Authors 	int		`json:"authors"`
	Fraction 	float32 `json:"fraction"`
}

type ReadershipCountry struct {
	Country 	string	`json:"country"`
	Iso 		string	`json:"isocode"`
	Books 		int		`json:"books"`
	Authors		int		`json:"authors"`
	Readership	int		`json:"readership"`
}

type StatusResponse struct{
	GutendexAPI		string	`json:"gutendexapi"`
	LanguageAPI		string	`json:"languageapi"`
	CountriesAPI	string	`json:"countriesapi"`
	Version 		string	`json:"version"`
	Uptime			string 	`json:"uptime"`

}

type RestCountriesResponse []struct {
	Population int `json:"population"`
}

type Country struct {
	Iso_three 		string  `json:"ISO3166_1_Alpha_3"`
	Iso_two			string  `json:"ISO3166_1_Alpha_2"`
	OfficialName	string  `json:"Official_Name"`
	Region			string	`json:"Region_Name"`
	SubRegion		string  `json:"Sub_Region_Name"`
	Language 		string	`json:"Language"`
}

type LanguageToCountriesResponse []Country

type Person struct {
	BirthYear 	int 	`json:"birth_year"`
	DeathYear 	int 	`json:"death_year"`
	Name 		string  `json:"name"`
}

type Book struct {
	Id    			int   				`json:"id"`
	Title			string				`json:"title"`
	Subjects		[]string			`json:"subjects"`
	Authors			[]Person			`json:"authors"`
	Translators		[]Person			`json:"translators"`
	Bookshelves		[]string			`json:"bookshelves"`
	Languages		[]string			`json:"languages"`
	Copyright		bool				`json:"copyright"`
	MediaType		map[string][]string	`json:"media_type"`
	Formats			string				`json:"formats"`
	DownloadCount	int					`json:"download_count"`
}

type GutenbergBookResponse struct {
	Count    	int       `json:"count"`
	Next    	string    `json:"next"`
	Previous    string    `json:"previous"`
	Books 		[]Book 	  `json:"results"`
}

//Functionality-----------------------------------------------------
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

func getTotalBooks() (int, error) {
	responseBooks, err := http.Get("http://129.241.150.113:8000/books")
	if err != nil {
		fmt.Print(err.Error())
		return 0, err
	}

	responseData, err := ioutil.ReadAll(responseBooks.Body)
	if err != nil {
		return 0, err
	}
	var response GutenbergBookResponse
	json.Unmarshal(responseData, &response)
	return response.Count, nil
}

func bookCount(w http.ResponseWriter, r *http.Request){
	var responseList []BookCountResponse
	var queryParameters = r.URL.Query()

	if !queryParameters.Has("language"){
		return
	}
	var totalBooks, err = getTotalBooks()
	if err != nil{
		returnErrorStatus(w, err.Error())
	}

	var languages = queryParameters.Get("language")
	var languageArray = strings.Split(languages, ",")


	for _, language := range languageArray{
		var response, err = bookCountForSingleLanguage(language, totalBooks)

		if err != nil {
			returnErrorStatus(w, err.Error())
			return
		}

		responseList = append(responseList, response)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseList)


}

func bookCountForSingleLanguage(language string, totalBooks int) (BookCountResponse, error){
	var page = 1
	var hasNextPage = true
	var allBooks []Book
	var count = 0

	for(hasNextPage){
		var url = fmt.Sprintf("http://129.241.150.113:8000/books/?languages=%s&page=%d",language,page)
		responseBooks, err:= http.Get(url)
		if err != nil{
			fmt.Print(err.Error())
			return BookCountResponse{}, err
		}

		responseData, err := ioutil.ReadAll(responseBooks.Body)
		if err != nil {
			return BookCountResponse{}, err
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

		return ourResponse, nil
}

func readership(w http.ResponseWriter, r *http.Request){
	var response []ReadershipCountry
	var limit = -1

	if (r.URL.Query().Has("limit")){
		limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	}

	var language = path.Base(r.URL.Path)
	var countries, err = getCountriesFromLanguage(language,limit)

	if err != nil{
		returnErrorStatus(w, err.Error())
	}



	//fraction will be wrong here but it is not needed
	var bookInfo, err2 = bookCountForSingleLanguage(language,100)

	if err != nil{
		returnErrorStatus(w, err2.Error())
	}

	for _, country := range countries{
		var population,error = getPopulationForCountry(country.OfficialName)

		if error != nil{
			returnErrorStatus(w,error.Error())
			return
		}

		var singleCountry = ReadershipCountry{
			country.OfficialName,
			country.Iso_two,
			bookInfo.Books,
			bookInfo.Authors,
			population,
		}

		response = append(response, singleCountry)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Println("Endpoint Hit: Readership")

}

func returnErrorStatus(w http.ResponseWriter,error string){
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w,error)
}

func getCountriesFromLanguage(language string, limit int) ([]Country, error){
	var url = fmt.Sprintf("http://129.241.150.113:3000/language2countries/%s",language)

	if (limit != -1){
		url = url + fmt.Sprintf("?limit=%d",limit)
	}


	var counter = 0
	var countries []Country
	responseCountries, err:= http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	responseData, err := ioutil.ReadAll(responseCountries.Body)
	if err != nil {
		return nil, err
	}


	var response LanguageToCountriesResponse
	json.Unmarshal(responseData, &response)


	for _, country := range response{
		countries = append(countries, country)
		counter++
		if ( limit != -1 && counter >= limit){
			return countries,nil
		}
	}

	return countries,nil

}

func getPopulationForCountry(countryName string) (int, error){
	var url = fmt.Sprintf("http://129.241.150.113:8080/v3.1/name/%s?fields=population",countryName)
	responseCountryData, err:= http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		return 0,err
	}


	responseData, err := ioutil.ReadAll(responseCountryData.Body)
	if err != nil {
		log.Fatal(err)
		return 0,err
	}

	var response RestCountriesResponse
	json.Unmarshal(responseData, &response)

	return response[0].Population,nil

}

func status(w http.ResponseWriter, r *http.Request){
	var gutenburgStatus = "0"
	var languagesToCountriesStatus = "0"
	var restCountriesStatus = "0"
	var version = "v1"

	//check availability of gutenberg
	responseBooks, err := http.Get("http://129.241.150.113:8000/books")
	if err != nil {
		fmt.Print(err.Error())
		gutenburgStatus = "500"
	} else {
		gutenburgStatus = responseBooks.Status
	}

	//check availability of languagetocountries
	responseLanguage, err := http.Get("http://129.241.150.113:3000/language2countries/no")
	if err != nil {
		fmt.Print(err.Error())
		languagesToCountriesStatus = "500"
	} else {
		languagesToCountriesStatus = responseLanguage.Status
	}

	//check availability of restCountries
	responseCountries, err := http.Get("http://129.241.150.113:8080/v3.1/name/Norway?fields=population")
	if err != nil {
		fmt.Print(err.Error())
		restCountriesStatus = "500"
	} else {
		restCountriesStatus = responseCountries.Status
	}


	var status = StatusResponse{
		gutenburgStatus,
		languagesToCountriesStatus,
		restCountriesStatus,
		version,
		"since March 6, 2024 at 11:18",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)

	fmt.Println("Endpoint Hit: Status")
}
//used for testing deployment on Render
func health(w http.ResponseWriter, r *http.Request){
}
//Handler for endpoints
func handleRequests() {
	http.HandleFunc("/librarystats/v1/bookcount/", bookCount)
	http.HandleFunc("/librarystats/v1/readership/", readership)
	http.HandleFunc("/librarystats/v1/status/", status)
	http.HandleFunc("/health/",health)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


func main() {
	handleRequests()
}
