package main

import (
	"Strava/tools"
	"net/http"
)

func handlerGetKudosRanking(w http.ResponseWriter, r *http.Request) {
	ranking, err := getKudosRanking()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tools.WriteJsonData(ranking, w)
}

func handlerStravaRequest(w http.ResponseWriter, r *http.Request, request string) {
	// Setup token in config.yaml to allow requests to Strava API
	token, err := setupConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := tools.ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var message string
	switch request {
	case "getActivity":
		go getActivity(db, token)
		message = "Activity processing from Strava API in progress"
	case "getKudos":
		go getKudos(db, token)
		message = "Kudos processing from Strava API in progress"
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest) // 400 Bad Request
		return
	}

	response := map[string]interface{}{
		"msg": message,
	}
	tools.WriteJsonData(response, w)
}

func handlerGetActivity(w http.ResponseWriter, r *http.Request) {
	handlerStravaRequest(w, r, "getActivity")
}

func handlerGetKudos(w http.ResponseWriter, r *http.Request) {
	handlerStravaRequest(w, r, "getKudos")
}
