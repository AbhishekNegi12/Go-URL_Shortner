package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// structuring the type
type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

// in memory database
// creating a map
// short-url will be mapped
var urlDB = make(map[string]URL)

// function to short url
// technique used is hahing algorithm
// this function will take url and give the short url
func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	// hasher.Write("hasher:",hasher)
	hasher.Write([]byte(OriginalURL)) //It converts the originalURl string to a byte slice
	// fmt.Println("hasher :", hasher)

	data := hasher.Sum(nil)
	// fmt.Println("hasher data :", data)

	//importing hasher data to string
	hash := hex.EncodeToString(data)
	// fmt.Println("EncodeToString data :", hash)

	//Final String
	// fmt.Println("Final String :", hash[:8])
	return hash[:8]
}

// after generate the short url should be stored in DB
func createURL(originalURl string) string {
	shortURL := generateShortURL(originalURl)

	//for the sake of simplicity we are making our id same as hashed url
	id := shortURL

	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURl,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

// getting the original URL corresponding to that shortURL
func getURL(id string)(URL, error){
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("url not found")
	}
	return url, nil
}

//	func handler(w http.ResponseWriter, r *http.Request){
//		fmt.Println("GET method")
//		// it write not in standard writer but the writer we provide
//		fmt.Fprintf(w, "Hello World!")
//	}
func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET method")
	// it write not in standard writer but the writer we provide
	fmt.Fprintf(w, "Hello World!")
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Welcome to the URL Shortener API"}
	json.NewEncoder(w).Encode(response)
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	shortURl_ := createURL(data.URL)
	// fmt.Fprintf(w, shorURl)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURl_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// redirecting URL Handler
func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
func main() {
	fmt.Println("Making of URL-Shortner")
	OriginalURL := "https://www.indiabix.com/"
	generateShortURL(OriginalURL)

	//Register the handler function to register all the request to the root url ("/")
	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	// Start the HTTP server
	fmt.Println("Starting the server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error on Starting server", err)
	}
}
