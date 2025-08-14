package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type GoogleMapsClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewGoogleMapsClient(apiKey string) *GoogleMapsClient {
	return &GoogleMapsClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *GoogleMapsClient) NearbySearch(lat, lng string, radius int, placeType string) (*PlacesResponse, error) {
	var urlBuilder strings.Builder
	urlBuilder.WriteString("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=")
	urlBuilder.WriteString(lat)
	urlBuilder.WriteString(",")
	urlBuilder.WriteString(lng)
	urlBuilder.WriteString("&radius=")
	urlBuilder.WriteString(strconv.Itoa(radius))
	urlBuilder.WriteString("&type=")
	urlBuilder.WriteString(placeType)
	urlBuilder.WriteString("&key=")
	urlBuilder.WriteString(c.apiKey)
	url := urlBuilder.String()

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// log.Println("Received response:", string(body))

	var placesResponse PlacesResponse
	if err := json.Unmarshal(body, &placesResponse); err != nil {
		return nil, err
	}

	return &placesResponse, nil
}

func (c *GoogleMapsClient) DistanceMatrix(originLat, originLng, placeID, mode, units string) (*DistanceMatrixResponse, error) {
	var distanceURLBuilder strings.Builder
	distanceURLBuilder.WriteString("https://maps.googleapis.com/maps/api/distancematrix/json?origins=")
	distanceURLBuilder.WriteString(originLat)
	distanceURLBuilder.WriteString(",")
	distanceURLBuilder.WriteString(originLng)
	distanceURLBuilder.WriteString("&destinations=place_id:")
	distanceURLBuilder.WriteString(placeID)
	distanceURLBuilder.WriteString("&mode=")
	distanceURLBuilder.WriteString(mode)
	distanceURLBuilder.WriteString("&units=")
	distanceURLBuilder.WriteString(units)
	distanceURLBuilder.WriteString("&key=")
	distanceURLBuilder.WriteString(c.apiKey)
	distanceURL := distanceURLBuilder.String()

	resp, err := c.httpClient.Get(distanceURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// log.Println("Received response:", string(body))

	var distanceResponse DistanceMatrixResponse
	if err := json.Unmarshal(body, &distanceResponse); err != nil {
		return nil, err
	}

	return &distanceResponse, nil
}
