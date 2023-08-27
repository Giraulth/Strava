package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type errorResponse struct {
	Message string              `json:"message"`
	Errors  []map[string]string `json:"errors"`
	Code    string              `json:"code"`
}

func CheckErrors(response []byte, clientId string) error {
	var errorResp errorResponse
	err := json.Unmarshal(response, &errorResp)
	if err != nil {
		return fmt.Errorf("Error parsing JSON response : %w", err)
	}
	if len(errorResp.Errors) > 0 {
		switch errorResp.Errors[0]["field"] {
		case "client_id":
			return fmt.Errorf("client_id is not valid in `config.yaml`. Go to https://www.strava.com/settings/api.")
		case "code":
			return fmt.Errorf("code is not valid in `config.yaml`. Go to https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=http://localhost/exchange_token&approval_prompt=force&scope=activity:read_all", clientId)
		default:
			return fmt.Errorf("Issue with field %s : %s", errorResp.Errors[0]["field"], errorResp.Errors[0]["code"])
		}
	}

	return nil
}

func SendRequest(method string, url string, requestBody []byte, token string) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return []byte{}, fmt.Errorf("Erreur creating request : %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("Error sending request : %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Error reading response : %w", err)
	}

	return body, nil
}

func WriteJsonData(data interface{}, w http.ResponseWriter) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // 200 OK
	// Set the Content-Type header to indicate JSON response
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "Internal Server Error "+err.Error(), http.StatusInternalServerError)
		return
	}
}
