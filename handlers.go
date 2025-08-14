package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type PubWithDistance struct {
	Place    `json:"-"`
	Name     string   `json:"name"`
	Vicinity string   `json:"vicinity"`
	Distance string   `json:"distance"`
	Duration string   `json:"duration"`
	Types    []string `json:"types"`
	PhotoURL string   `json:"photo_url"`
}

func getPlaces(client *GoogleMapsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		lng := r.URL.Query().Get("lng")

		if lat == "" || lng == "" {
			http.Error(w, "lat and lng are required", http.StatusBadRequest)
			return
		}

		placesResponse, err := client.NearbySearch(lat, lng, 1609, "bar")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var pubsWithDistance []PubWithDistance

		for _, place := range placesResponse.Results {
			distanceResponse, err := client.DistanceMatrix(lat, lng, place.PlaceID, "walking", "imperial")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(distanceResponse.Rows) > 0 && len(distanceResponse.Rows[0].Elements) > 0 {
				var photoURL string
				if len(place.Photos) > 0 {
					var photoURLBuilder strings.Builder
					photoURLBuilder.WriteString("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=")
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
					Types:    place.Types,
					PhotoURL: photoURL,
				})
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pubsWithDistance)
	}
}