package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type PubWithDistance struct {
	Place    `json:"-"`
	Name     string `json:"name"`
	Vicinity string `json:"vicinity"`
	Distance string `json:"distance"`
	Duration string `json:"duration"`
	PhotoURL string `json:"photo_url"`
}

func getPlacesInternal(client *GoogleMapsClient, lat, lng string) ([]PubWithDistance, error) {
	placesResponse, err := client.NearbySearch(lat, lng, "bar")
	if err != nil {
		return nil, err
	}

	var pubsWithDistance []PubWithDistance

	for _, place := range placesResponse.Results {
		distanceResponse, err := client.DistanceMatrix(lat, lng, place.PlaceID, "walking", "imperial")
		if err != nil {
			return nil, err
		}

		if len(distanceResponse.Rows) > 0 && len(distanceResponse.Rows[0].Elements) > 0 {
			var photoURL string
			if len(place.Photos) > 0 {
				var photoURLBuilder strings.Builder
				photoURLBuilder.WriteString(GooglePlacePhotoURL)
				photoURLBuilder.WriteString("?maxwidth=400&photoreference=")
				photoURLBuilder.WriteString(place.Photos[0].PhotoReference)
				photoURLBuilder.WriteString("&key=")
				photoURLBuilder.WriteString(client.apiKey)
				photoURL = photoURLBuilder.String()
			}

			pubsWithDistance = append(pubsWithDistance, PubWithDistance{
				Place:    place,
				Name:     place.Name,
				Vicinity: place.Vicinity,
				Distance: distanceResponse.Rows[0].Elements[0].Distance.Text,
				Duration: distanceResponse.Rows[0].Elements[0].Duration.Text,
				PhotoURL: photoURL,
			})
		}
	}
	return pubsWithDistance, nil
}

func getPlaces(client *GoogleMapsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		lng := r.URL.Query().Get("lng")

		if lat == "" || lng == "" {
			http.Error(w, "lat and lng are required", http.StatusBadRequest)
			return
		}

		pubs, err := getPlacesInternal(client, lat, lng)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pubs)
	}
}

func getLocationFromIP(ip string) (*Location, error) {
	// For testing purposes, since localhost will resolve to ::1 or 127.0.0.1
	if ip == "::1" || ip == "127.0.0.1" {
		ip = "35.242.175.203" // GOOGLE-CLOUD-PLATFORM London DNS
		log.Printf("Using test IP address: %s", ip)
	}

	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// log.Printf("ip-api.com response: %s", string(body))

	var location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	if err := json.Unmarshal(body, &location); err != nil {
		return nil, err
	}

	return &Location{Lat: location.Lat, Lng: location.Lon}, nil
}

func getLocation(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	log.Printf("IP address: %s", ip)

	location, err := getLocationFromIP(ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}

func getPlacesNearMe(client *GoogleMapsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}

		log.Printf("IP address: %s", ip)

		location, err := getLocationFromIP(ip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lat := strconv.FormatFloat(location.Lat, 'f', -1, 64)
		lng := strconv.FormatFloat(location.Lng, 'f', -1, 64)

		pubs, err := getPlacesInternal(client, lat, lng)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pubs)
	}
}
