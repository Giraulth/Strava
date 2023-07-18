package main

import (
	"Strava/tools"
	"encoding/json"
	"fmt"
	"net/url"
)

var stravaApi string = "https://www.strava.com/api/v3/"

func getActivity(token string) error {
	baseUrl := stravaApi + "/athlete/activities?"

    for i := 1; i <= 20; i++ {
		queryParams := url.Values{
			"page": {fmt.Sprintf("%d", i)},
		}

		requestUrl := baseUrl + queryParams.Encode()
		response, err := tools.SendRequest("GET", requestUrl, nil, token)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	
		var activities []map[string]interface{}
		err = json.Unmarshal(response, &activities)
		if err != nil {
			err = tools.CheckErrors(response, "")
			if err != nil {
				return fmt.Errorf("%v", err)
			}
		}
	
		if len(activities) == 0 {
			return nil // all activities have been retrieved
		}
		for _, activity := range activities {
			fmt.Printf("%.0f\n", activity["id"].(float64))
		}
    }


	return nil
}
