package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type SignupPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (app *Config) Signup(w http.ResponseWriter, r *http.Request) {

	//extract the request body
	var requestPayload SignupPayload

	//extract the requestbody
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, nil)
		return
	}

	// Validate the request payload
	if err := app.ValidataSignupInput(requestPayload); len(err) > 0 {
		log.Println(err)
		app.errorJSON(w, errors.New("error trying to sign-up user"), err, http.StatusBadRequest)
		return
	}

	//create some json we will send to authservice
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")

	authServiceUrl := fmt.Sprintf("%s%s", os.Getenv("AUTH_URL"), "signup")

	log.Println(authServiceUrl)

	// call the service by creating a request
	request, err := http.NewRequest("POST", authServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, nil)
		return
	}

	// Set the Content-Type header
	request.Header.Set("Content-Type", "application/json")
	//create a http client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, nil)
		return
	}
	defer response.Body.Close()

	// create a varabiel we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err, nil)
		return
	}

	log.Println("response from auth service", jsonFromService)
	if response.StatusCode != http.StatusAccepted {
		log.Println(jsonFromService.Message, jsonFromService)
		app.errorJSON(w, errors.New(jsonFromService.Message), nil)
		return
	}

	var payload jsonResponse
	payload.Error = jsonFromService.Error
	payload.StatusCode = http.StatusOK
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {

	//extract the requestbody
	var requestPayload LoginPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, nil)
		return
	}

	//create some json we will send to authservice
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")

	//call the service
	// call the service by creating a request
	request, err := http.NewRequest("POST", os.Getenv("AUTH_URL")+"login", bytes.NewBuffer(jsonData))

	log.Println("reached the handler in broker for login")
	if err != nil {
		app.errorJSON(w, err, nil)
		return
	}

	// Set the Content-Type header
	request.Header.Set("Content-Type", "application/json")
	//create a http client
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err, nil)
		return
	}
	defer response.Body.Close()

	// create a variable that 'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err, nil)
		return
	}

	//check the status of the response
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New(jsonFromService.Message), nil)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, nil, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = jsonFromService.Error
	payload.StatusCode = http.StatusOK
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	result, err := app.getUserToken(w, r)
	if err != nil {
		app.errorJSON(w, errors.New(result.Message), nil)
		return
	}
	if result.Error {
		app.errorJSON(w, errors.New(result.Message), result.Data)
		return
	}
	// call the service by creating a request
	request, err := http.NewRequest("GET", os.Getenv("AUTH_URL")+"all-users", nil)

	if err != nil {
		app.errorJSON(w, err, nil)
		return
	}

	// Set the Content-Type header
	request.Header.Set("Content-Type", "application/json")
	//create a http client
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err, nil)
		return
	}
	defer response.Body.Close()

	// create a variable that 'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err, nil)
		return
	}

	//check the status of the response
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New(jsonFromService.Message), nil)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, nil, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = jsonFromService.Error
	payload.StatusCode = http.StatusOK
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)
}

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
func (app *Config) GetMe()       {}
func (app *Config) VerifyToken() {}
func (app *Config) Logout()      {}
