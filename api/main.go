package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/strava/getActivity", handlerGetActivity)
	http.HandleFunc("/strava/getKudos", handlerGetKudos)
	http.HandleFunc("/api/getKudosRanking", handlerGetKudosRanking)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)

}
