package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type JSONResponse struct {
	Message string `json:"message"`
}

func handlerStravaRequest(w http.ResponseWriter, r *http.Request, request string) {
	// Setup token in config.yaml to allow requests to Strava API
	token, err := setupConfig()
	if err != nil {
		fmt.Println(err)
	}

	message := "Success"
	switch request {
	case "getActivity":
		go getActivity(token)
		message = "Activity processing from Strava API in progress"
	case "getKudos":
		go getKudos(token)
		message = "Kudos processing from Strava API in progress"
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest) // 400 Bad Request
		return
	}

	w.WriteHeader(http.StatusOK) // 200 OK
	// Create the JSON response
	response := JSONResponse{
		Message: message,
	}
	// Convert the response struct to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to indicate JSON response
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func HandlerGetActivity(w http.ResponseWriter, r *http.Request) {
	handlerStravaRequest(w, r, "getActivity")
}

func HandlerGetKudos(w http.ResponseWriter, r *http.Request) {
	handlerStravaRequest(w, r, "getKudos")
}
