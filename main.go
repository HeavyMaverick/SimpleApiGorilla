package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Movie struct {
	Title  string `json:"title,omitempty"`
	Rating int    `json:"rating,omitempty"`
}

var (
	mu     sync.RWMutex
	movies []Movie = []Movie{
		{"Marti Supreme", 10},
		{"War Dogs", 10},
	}
)

func GetMovies(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		log.Println(err)
		return
	}
}
func GetMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mu.RLock()
	defer mu.RUnlock()
	for _, item := range movies {
		if item.Title == params["title"] {
			err := json.NewEncoder(w).Encode(&item)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
}
func AddMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	movies = append(movies, movie)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		log.Println(err)
		return
	}
}
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mu.Lock()
	defer mu.Unlock()
	for index, item := range movies {
		if item.Title == params["title"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		log.Println(err)
		return
	}

}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Movies API: GET /movies, GET /movies/{title}, POST /movies, DELETE /movies/{title}"))
		if err != nil {
			log.Println(err)
		}
	})
	router.HandleFunc("/movies", GetMovies).Methods("GET")
	router.HandleFunc("/movies/{title}", GetMovie).Methods("GET")
	router.HandleFunc("/movies", AddMovie).Methods("POST")
	router.HandleFunc("/movies/{title}", DeleteMovie).Methods("DELETE")
	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
