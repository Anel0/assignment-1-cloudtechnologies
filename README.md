# Assignment 1



## Author - Ane Lo



## Description

Welcome to the Book Information API documentation! This REST web application, written in Golang, serves as a gateway to access specific information about books from the Gutendex library. The service is built on the data from three pre-existing APIs:


- Gutendex API
    - Endpoint: http://129.241.150.113:8000/books/
    - Documentation: http://129.241.150.113:8000/


- Language2Countries API
    - Endpoint: http://129.241.150.113:3000/language2countries/
    - Documentation: http://129.241.150.113:3000/


- REST Countries API
    - Endpoint: http://129.241.150.113:8080/v3.1
    - Documentation: http://129.241.150.113:8080/


The application has three endpoints of its own, listed and described below.


# Endpoints



## Bookcount Endpoint

Like in the assignment description, this initial endpoint returns the number of books for any given language, identified via country 2-letter language ISO codes (ISO 639 Set 1), as well as the number of unique authors. This can be a single as well as multiple languages (comma-separated language codes).


**Request:**

```
Method: GET
Path: bookcount/?language={:two_letter_language_code+}/
```



two_letter_language_code is the corresponding 2-letter language ISO codes (ISO 639 Set 1)


**Example requests:**

```
bookcount/?language=no
bookcount/?language=no,fi
```



**Response**

Content type: application/json

Body (Example):

```
[
  {
     "language": "no",
     "books": 21,
     "authors": 14,
     "fraction": 0.0005
  },
  {
     "language": "fi",
     "books": 2798,
     "authors": 228,
     "fraction": 0.0671
  }
]
```



checklist:

- [ ] The language code is the same as the input code.
- [ ] The book and author information is highlighting the number of available books and authors (from gutendex) in the queried language.
- [ ] The fraction is the number of books divided by all books served via gutendex.
- [ ] The service should be generic (i.e., allow open entry for any language)
- [ ] As indicated above, the service should allow for single or multiple languages.



## Readership Endpoint

This second endpoint returns the number of potential readers for books in a given language, i.e., the population per country in which that language is official (and hence assuming that the inhabitants can potentially read it).the number of potential readers is provided together with the number of books and authors associated with a given language. 

**Request**

```
Method: GET
Path: readership/{:two_letter_language_code}{?limit={:number}}
```


- {:two_letter_language_code} refers to the ISO639 Set 1 identifier of the language for which you establish readership.

- {?limit={:number}} is an optional parameter that limits the number of country entries that are reported (in addition to the total number).

**Example requests:**

```
readership/no
readership/no/?limit=5
```


**Response**

Content type: application/json

Status code: gives 200 if everything is OK


Body (Example):

```
[ 
  {
     "country": "Norway",
     "isocode": "NO",
     "books": 21,
     "authors": 14,
     "readership": 5379475
  },
  {
     "country": "Svalbard and Jan Mayen",
     "isocode": "SJ",
     "books": 21,
     "authors": 14,
     "readership": 2562
  },
  {
     "country": "Iceland",
     "isocode": "IS",
     "books": 21,
     "authors": 14,
     "readership": 366425
  }
]
```


checklist:

- [ ] The isocode field should be the two-letter ISO code for the country (recall: country != language) (as per ISO3166-1-Alpha-2).
- [ ] The readership is the population of that country.
- [ ] The country name should be the English country name.
- [ ] The book and author information is - as with the previous endpoint - highlighting the number of available books and authors (from gutendex) in the queried language.
- [ ] This endpoint only needs to support input of a single (not multiple) languages at a given time.



## Status Endpoint
The diagnostics interface indicates the availability of individual services this service depends on. The reporting occurs based on status codes returned by the dependent services, and it further provides information about the uptime of the service.

**Request**

```
Method: GET
Path: status/
```


**Response**

Content type: application/json

Status code: 200 if everything is OK

Body:

```
{
   "gutendexapi": "<http status code for gutendex API>",
   "languageapi: "<http status code for language2countries API>", 
   "countriesapi": "<http status code for restcountries API>",
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}
```

Note: < some value > indicates placeholders for values to be populated by the service.



## How to use


The application is reployed on Render 
- URL: https://assignment-1-cloudtechnologies.onrender.com

Exmples of how to use web service:

https://assignment-1-cloudtechnologies.onrender.com/librarystats/v1/status/
https://assignment-1-cloudtechnologies.onrender.com/librarystats/v1/bookcount/?language=no,ko,se
https://assignment-1-cloudtechnologies.onrender.com/librarystats/v1/readership/se?limit=1



Localhost:

to test the program on your own, clone the git repo,
open the program in your preferred IDE and run the program. 
The default server runs on https://localhost8080
 
```
http://localhost:8080/librarystats/v1/bookcount/
http://localhost:8080/librarystats/v1/readership/
http://localhost:8080/librarystats/v1/status/
```





## Support
If you have questions or feedback, please email Ane at:
- anlo@stud.ntnu.no

