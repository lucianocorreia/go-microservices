package main

import (
	"logger/data"
	"net/http"
)

type logRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json
	var req logRequest
	_ = app.readJSON(w, r, &req)

	// insert
	event := data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
