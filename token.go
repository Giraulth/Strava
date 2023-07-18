package main

import (
	"Strava/tools"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var stravaUrl string = "https://www.strava.com/oauth/token?"
var configFilename string = "config.yaml"

type requestBody struct {
	User struct {
		ClientId     int    `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		Code         string `yaml:"code"`
		AccessToken  string `yaml:"access_token"`
		RefreshToken string `yaml:"refresh_token"`
	} `yaml:"user"`
}

func updateConfig(requestBody requestBody, access_token string, refresh_token string) error {

	requestBody.User.AccessToken = access_token
	requestBody.User.RefreshToken = refresh_token

	// Marshal the updated Config struct into YAML
	updatedYAML, err := yaml.Marshal(&requestBody)
	if err != nil {
		return fmt.Errorf("Failed to marshal YAML: %v", err)
	}

	// Write the updated YAML content back to the file
	err = ioutil.WriteFile(configFilename, updatedYAML, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to write YAML: %v", err)
	}

	return nil
}

func getToken(requestBody requestBody, grantType string) (string, error) {

	clientId := fmt.Sprintf("%d", requestBody.User.ClientId)
	clientSecret := requestBody.User.ClientSecret
	code := requestBody.User.Code
	refreshToken := requestBody.User.RefreshToken

	queryParams := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"grant_type":    {grantType},
	}

	if grantType == "refresh_token" {
		queryParams.Set("refresh_token", refreshToken)
	} else {
		queryParams.Set("code", code)
	}
	tokenUrl := stravaUrl + queryParams.Encode()

	response, err := tools.SendRequest(http.MethodPost, tokenUrl, nil, "")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	var token map[string]interface{}
	err = json.Unmarshal(response, &token)
	if err != nil {
		return "", fmt.Errorf("Error parsing JSON response : %w", err)
	}

	err = tools.CheckErrors(response, clientId)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	accessToken, ok := token["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token is not a string")
	}

	refreshToken, ok = token["refresh_token"].(string)
	if !ok {
		return "", fmt.Errorf("refresh_token is not a string")
	}

	// todo : manage errors
	err = updateConfig(requestBody, accessToken, refreshToken)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return accessToken, nil
}

func setupConfig() (string, error) {

	var requestBody requestBody
	err := tools.ReadConfig(configFilename, &requestBody)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	access_token := requestBody.User.AccessToken
	refresh_token := requestBody.User.RefreshToken

	if access_token == "" && refresh_token == "" {
		_, err = getToken(requestBody, "authorization_code")
		if err != nil {
			return "", fmt.Errorf("%v", err)
		}
	}

	token, err := getToken(requestBody, "refresh_token")
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return token, nil
}
