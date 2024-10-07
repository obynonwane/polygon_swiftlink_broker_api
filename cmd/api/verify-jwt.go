package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func (app *Config) getUserToken(w http.ResponseWriter, r *http.Request) (jsonResponse, error) {
	//get authorization hearder
	authorizationHeader := r.Header.Get("Authorization")

	// call the service by creating a request
	request, err := http.NewRequest("GET", os.Getenv("AUTH_URL")+"verify-user-token", nil)

	if err != nil {
		return jsonResponse{Error: true, Message: err.Error(), StatusCode: http.StatusBadRequest, Data: nil}, err

	}

	// Set the "Authorization" header with your Bearer token
	request.Header.Set("authorization", authorizationHeader)

	// Set the Content-Type header
	request.Header.Set("Content-Type", "application/json")
	//create a http client
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return jsonResponse{Error: true, Message: err.Error(), StatusCode: http.StatusBadRequest, Data: nil}, err

	}
	defer response.Body.Close()

	//variable to marshal into
	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		return jsonResponse{Error: true, Message: err.Error(), StatusCode: http.StatusBadRequest, Data: nil}, err
	}

	// make a call to the bank-service
	var payload jsonResponse
	payload.Error = jsonFromService.Error
	payload.Message = jsonFromService.Message
	payload.StatusCode = response.StatusCode
	payload.Data = jsonFromService.Data

	if jsonFromService.Error {
		return payload, err
	}

	return payload, nil
}
