package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	apiKey         = "AIzaSyDUeHxBudfbcj_04PAV59MEyS5UWHKmK6I"
	radiusInMeters = 1609 // 1 mile in meters
)

func main() {
	r := mux.NewRouter()
	client := NewGoogleMapsClient(apiKey, radiusInMeters)
	r.HandleFunc("/places", getPlaces(client))
	r.HandleFunc("/location", getLocation)
	r.HandleFunc("/places/near-me", getPlacesNearMe(client))
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
