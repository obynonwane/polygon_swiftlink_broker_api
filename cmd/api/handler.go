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

func (app *Config) GetMe()       {}
func (app *Config) VerifyToken() {}
func (app *Config) Logout()      {}

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
func (app *Config) MainnetMissedCheckpoint(w http.ResponseWriter, r *http.Request) {

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
	request, err := http.NewRequest("GET", os.Getenv("SERVICE_URL")+"pos/mainnet/mainnet-missed-checkpoint", nil)

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
func (app *Config) TestnetMissedCheckpoint(w http.ResponseWriter, r *http.Request) {

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
	request, err := http.NewRequest("GET", os.Getenv("SERVICE_URL")+"pos/testnet/mainnet-missed-checkpoint", nil)

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
func (app *Config) MainnetHeimdalBlockHeight(w http.ResponseWriter, r *http.Request) {

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
	request, err := http.NewRequest("GET", os.Getenv("SERVICE_URL")+"pos/mainnet/heimdal-block-height", nil)

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
func (app *Config) TestnetHeimdalBlockHeight(w http.ResponseWriter, r *http.Request) {

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
	request, err := http.NewRequest("GET", os.Getenv("SERVICE_URL")+"pos/testnet/heimdal-block-height", nil)

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
func (app *Config) MainnetBorLatestBlockDetails(w http.ResponseWriter, r *http.Request) {

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
	request, err := http.NewRequest("GET", os.Getenv("SERVICE_URL")+"pos/testnet/bor-latest-block-details", nil)

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
func (app *Config) TestnetBorLatestBlockDetails(w http.ResponseWriter, r *http.Request) {

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
	request, err := http.NewRequest("GET", os.Getenv("SERVICE_URL")+"pos/testnet/bor-latest-block-details", nil)

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
func (app *Config) MainnetStateSync(w http.ResponseWriter, r *http.Request) {

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
	request, err := http.NewRequest("GET", os.Getenv("SERVICE_URL")+"pos/mainnet/state-sync", nil)

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
func (app *Config) TestnetStateSync(w http.ResponseWriter, r *http.Request) {

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
	request, err := http.NewRequest("GET", os.Getenv("SERVICE_URL")+"pos/testnet/state-sync", nil)

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
