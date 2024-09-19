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
func (app *Config) Login()       {}
func (app *Config) GetMe()       {}
func (app *Config) VerifyToken() {}
func (app *Config) Logout()      {}
