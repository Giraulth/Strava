package main

import (
	"Strava/tools"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"net/url"
)

var stravaApi string = "https://www.strava.com/api/v3/"

func getDefaultFloatValue(activity map[string]interface{}, columnName string) float64 {
	if floatValue, ok := activity[columnName].(float64); ok {
		return floatValue
	}
	return 0.0
}

func getKudos(token string) error {
	baseUrl := stravaApi + "activities/"
	db := tools.ConnectDB()
	tableColumns := []string{"activity_id", "username"}

	activities := tools.SelectId(db)
	for _, id := range activities {
		requestUrl := baseUrl +  strconv.FormatInt(id, 10) + "/kudos"
		response, err := tools.SendRequest("GET", requestUrl, nil, token)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		var kudosList []map[string]interface{}
		err = json.Unmarshal(response, &kudosList)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		fmt.Printf("Processing %d : Received %d kudos\n", id, len(kudosList))
		for _, kudos := range kudosList {
			personName := fmt.Sprintf("%s %s", kudos["firstname"].(string), kudos["lastname"].(string))
			values := []interface{}{id, personName}
			err = tools.GenericInsert(db, "kudos", tableColumns, values...)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
		}
		if len(kudosList) == 0 {
			err = tools.GenericInsert(db, "kudos", tableColumns, []interface{}{id, ""}...)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}

func getActivity(token string) error {
	baseUrl := stravaApi + "/athlete/activities?"
	db := tools.ConnectDB()

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

			activityId := int64(activity["id"].(float64))
			activityName := activity["name"].(string)
			tableColumns := []string{"start_date", "id", "name", "type", "distance", "moving_time",
				"elapsed_time", "average_heartrate", "max_heartrate", "got_kudos_list"}
			values := []interface{}{activity["start_date_local"].(string),
				activityId,
				activityName,
				activity["type"].(string),
				activity["distance"].(float64),
				activity["moving_time"].(float64),
				activity["elapsed_time"].(float64),
				getDefaultFloatValue(activity, "average_heartrate"),
				getDefaultFloatValue(activity, "max_heartrate"),
				false}
			fmt.Printf("Activity in progress: %s (%d) \n", activityName, activityId)
			err := tools.GenericInsert(db, "activity", tableColumns, values...)
			if err != nil {
				return fmt.Errorf("Activity %d is already in database", activityId)
			}

		}
	}

	defer db.Close()
	return nil
}
