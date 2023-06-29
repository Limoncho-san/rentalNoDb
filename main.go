package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Rental struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Type            string  `json:"type"`
	Make            string  `json:"make"`
	Model           string  `json:"model"`
	Year            int     `json:"year"`
	Length          float64 `json:"length"`
	Sleeps          int     `json:"sleeps"`
	PrimaryImageURL string  `json:"primary_image_url"`
	Price           struct {
		Day int `json:"day"`
	} `json:"price"`
	Location struct {
		City    string  `json:"city"`
		State   string  `json:"state"`
		Zip     string  `json:"zip"`
		Country string  `json:"country"`
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
	} `json:"location"`
	User struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"user"`
}

var rentals []Rental

func main() {
	// Initialize the rentals slice (replace with your actual data source)
	rentals = []Rental{
		{
			ID:              1,
			Name:            "Rental 1",
			Description:     "Description of Rental 1",
			Type:            "test1",
			Make:            "",
			Model:           "",
			Year:            0,
			Length:          0,
			Sleeps:          0,
			PrimaryImageURL: "",
			Price: struct {
				Day int `json:"day"`
			}{},
			Location: struct {
				City    string  `json:"city"`
				State   string  `json:"state"`
				Zip     string  `json:"zip"`
				Country string  `json:"country"`
				Lat     float64 `json:"lat"`
				Lng     float64 `json:"lng"`
			}{},
			User: struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			}{},
		},
		{
			ID:              2,
			Name:            "Rental 2",
			Description:     "Description of Rental 2",
			Type:            "",
			Make:            "",
			Model:           "",
			Year:            0,
			Length:          0,
			Sleeps:          0,
			PrimaryImageURL: "",
			Price: struct {
				Day int `json:"day"`
			}{
				Day: 0,
			},
			Location: struct {
				City    string  `json:"city"`
				State   string  `json:"state"`
				Zip     string  `json:"zip"`
				Country string  `json:"country"`
				Lat     float64 `json:"lat"`
				Lng     float64 `json:"lng"`
			}{
				City:    "",
				State:   "",
				Zip:     "",
				Country: "",
				Lat:     0,
				Lng:     0,
			},
			User: struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			}{
				ID:        0,
				FirstName: "",
				LastName:  "",
			},
		},
	}

	http.HandleFunc("/rentals/", handleGetRental)
	http.HandleFunc("/rentals", handleListRentals)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetRental(w http.ResponseWriter, r *http.Request) {
	rentalID := strings.TrimPrefix(r.URL.Path, "/rentals/")
	id, err := strconv.Atoi(rentalID)
	if err != nil {
		http.Error(w, "Invalid rental ID", http.StatusBadRequest)
		return
	}

	rental := findRentalByID(id)
	if rental == nil {
		http.Error(w, "Rental not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rental)
}

func handleListRentals(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters
	query := r.URL.Query()
	priceMin, _ := strconv.Atoi(query.Get("price_min"))
	priceMax, _ := strconv.Atoi(query.Get("price_max"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	ids := parseCommaSeparatedList(query.Get("ids"))
	near := parseCommaSeparatedPair(query.Get("near"))
	sort := query.Get("sort")

	// Filter rentals based on query parameters
	filteredRentals := filterRentals(priceMin, priceMax, ids, near)

	// Sort rentals if 'sort' parameter is provided
	if sort != "" {
		sortRentals(filteredRentals, sort)
	}

	// Apply limit and offset
	paginatedRentals := paginateRentals(filteredRentals, limit, offset)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedRentals)
}

func findRentalByID(id int) *Rental {
	for _, rental := range rentals {
		if rental.ID == id {
			return &rental
		}
	}
	return nil
}

func filterRentals(priceMin, priceMax int, ids []int, near []float64) []Rental {
	var filtered []Rental

	for _, rental := range rentals {
		if priceMin > 0 && rental.Price.Day < priceMin {
			continue
		}
		if priceMax > 0 && rental.Price.Day > priceMax {
			continue
		}
		if len(ids) > 0 && !containsID(ids, rental.ID) {
			continue
		}
		if len(near) == 2 && !isNear(rental.Location.Lat, rental.Location.Lng, near[0], near[1]) {
			continue
		}

		filtered = append(filtered, rental)
	}

	return filtered
}

func sortRentals(rentals []Rental, sort string) {
	// Implement sorting logic based on the 'sort' parameter
	// This example assumes sorting by price in ascending order
}

func paginateRentals(rentals []Rental, limit, offset int) []Rental {
	if limit <= 0 {
		return rentals
	}

	start := offset
	end := offset + limit
	if start >= len(rentals) {
		return []Rental{}
	}
	if end > len(rentals) {
		end = len(rentals)
	}

	return rentals[start:end]
}

func parseCommaSeparatedList(input string) []int {
	var result []int
	if input == "" {
		return result
	}
	values := strings.Split(input, ",")
	for _, value := range values {
		id, err := strconv.Atoi(value)
		if err == nil {
			result = append(result, id)
		}
	}
	return result
}

func parseCommaSeparatedPair(input string) []float64 {
	var result []float64
	if input == "" {
		return result
	}
	values := strings.Split(input, ",")
	for _, value := range values {
		num, err := strconv.ParseFloat(value, 64)
		if err == nil {
			result = append(result, num)
		}
	}
	return result
}

func containsID(ids []int, id int) bool {
	for _, value := range ids {
		if value == id {
			return true
		}
	}
	return false
}

func isNear(lat1, lng1, lat2, lng2 float64) bool {
	// Implement distance calculation logic here to determine if two locations are near
	return true
}