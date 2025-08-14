package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	apiKey = ""
)

func main() {
	r := mux.NewRouter()
	client := NewGoogleMapsClient(apiKey)
	r.HandleFunc("/places", getPlaces(client))
	r.HandleFunc("/location", getLocation)
	r.HandleFunc("/places/near-me", getPlacesNearMe(client))
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
