package main

import (
	"fmt"
)

func main() {
	// Setup token in config.yaml to allow requests to Strava API
	token, err := setupConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = getActivity(token)
	if err != nil {
		fmt.Println(err)
		return
	}
}
