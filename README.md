# FindMyPub

FindMyPub is a simple Go application that helps you find nearby pubs. It uses the Google Maps API to find places and provide distance and duration information.

## Getting Started

To get started, you will need to have Go installed on your machine.

1.  **Get a Google Maps API Key:**
    *   Go to the [Google Cloud Console](https://console.cloud.google.com/).
    *   Create a new project or select an existing one.
    *   Enable the "Places API" and "Distance Matrix API".
    *   Create an API key.
2.  **Update the API Key:**
    *   Open the `main.go` file.
    *   Replace the placeholder `apiKey` with your actual Google Maps API key.
3.  **Run the application:**
    ```bash
    go run .
    ```
    The server will start on `localhost:8080`.

## API Endpoints

### GET /location

Returns an estimated location (latitude and longitude) based on the user's IP address.

**Example Response:**

```json
{
  "lat": 37.422,
  "lng": -122.084
}
```

### GET /places

Returns a list of nearby pubs based on the provided latitude and longitude.

**Query Parameters:**

*   `lat` (required): Latitude
*   `lng` (required): Longitude

**Example Request:**

```
GET /places?lat=37.422&lng=-122.084
```

**Example Response:**

```json
[
  {
    "name": "The Old Pro",
    "vicinity": "541 Ramona St, Palo Alto",
    "distance": "0.5 mi",
    "duration": "10 mins",
    "types": [
      "bar",
      "restaurant",
      "food",
      "point_of_interest",
      "establishment"
    ],
    "photo_url": "https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=..."
  }
]
```

### GET /places/near-me

Returns a list of nearby pubs based on the user's estimated location (from their IP address). This endpoint combines the functionality of `/location` and `/places`.

**Example Request:**

```
GET /places/near-me
```

**Example Response:**

The response will be the same as the `/places` endpoint.
