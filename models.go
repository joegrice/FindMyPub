package main

type PlacesResponse struct {
	Results []Place `json:"results"`
}

type Place struct {
	Name      string   `json:"name"`
	Vicinity  string   `json:"vicinity"`
	Geometry  Geometry `json:"geometry"`
	PlaceID   string   `json:"place_id"`
	Types     []string `json:"types"`
	Photos    []Photo  `json:"photos"`
}

type Photo struct {
	PhotoReference string `json:"photo_reference"`
}

type Geometry struct {
	Location Location `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type DistanceMatrixResponse struct {
	Rows []struct {
		Elements []struct {
			Distance struct {
				Text string `json:"text"`
			} `json:"distance"`
			Duration struct {
				Text string `json:"text"`
			} `json:"duration"`
		} `json:"elements"`
	} `json:"rows"`
}
