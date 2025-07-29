package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type LocationIndex struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type LocationsData struct {
	Index []LocationIndex `json:"index"`
}

type DateIndex struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesData struct {
	Index []DateIndex `json:"index"`
}

type RelationIndex struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationsData struct {
	Index []RelationIndex `json:"index"`
}

type PageData struct {
	Artists []Artist
}

func fetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	return artists, err
}

func fetchLocations() (map[int][]string, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data LocationsData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	locations := make(map[int][]string)
	for _, loc := range data.Index {
		locations[loc.ID-1] = loc.Locations
	}
	return locations, nil
}

func fetchDates() (map[int][]string, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data DatesData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	dates := make(map[int][]string)
	for _, d := range data.Index {
		dates[d.ID-1] = d.Dates
	}
	return dates, nil
}

func fetchRelations() (map[int]map[string][]string, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data RelationsData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	relations := make(map[int]map[string][]string)
	for _, r := range data.Index {
		relations[r.ID-1] = r.DatesLocations
	}
	return relations, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := fetchArtists()
	if err != nil {
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, PageData{Artists: artists})
}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/locations/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	locations, err := fetchLocations()
	if err != nil {
		http.Error(w, "Failed to load locations", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Locations []string `json:"locations"`
	}{
		Locations: locations[id],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func datesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/dates/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	dates, err := fetchDates()
	if err != nil {
		http.Error(w, "Failed to load dates", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Dates []string `json:"dates"`
	}{
		Dates: dates[id],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func relationHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/relation/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	relations, err := fetchRelations()
	if err != nil {
		http.Error(w, "Failed to load relations", http.StatusInternalServerError)
		return
	}

	resp := struct {
		DatesLocations map[string][]string `json:"datesLocations"`
	}{
		DatesLocations: relations[id],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))


	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/locations/", locationsHandler)
	http.HandleFunc("/dates/", datesHandler)
	http.HandleFunc("/relation/", relationHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
