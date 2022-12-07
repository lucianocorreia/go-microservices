package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPauload RequestPayload

	err := app.readJSON(w, r, &requestPauload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPauload.Action {
	case "auth":
		app.authenticate(w, requestPauload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, payload AuthPayload) {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New(fmt.Sprintf("error calling auth service, %v", response.StatusCode)))
		return
	}

	var resp jsonResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if resp.Error {
		app.errorJSON(w, errors.New(resp.Message))
		return
	}

	var jsonPayload jsonResponse
	jsonPayload.Error = false
	jsonPayload.Message = "Authenticated!"
	jsonPayload.Data = resp.Data

	app.writeJSON(w, http.StatusAccepted, jsonPayload)
}
