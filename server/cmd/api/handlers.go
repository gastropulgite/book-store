package main

import (
	"net/http"
)

type cridentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type jsonResponse struct {
	Data interface{} `json:"data,omitempty"`
	Error bool `json:"error"`
	Message string `json:"message"`
}

type envelope map[string]interface{}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	
	var creds cridentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds) 

	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "Invalid JSON"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	app.infoLog.Println(creds.UserName, creds.Password)

	payload.Error = false
	payload.Message = "Signed In!"

	err = app.writeJSON(w, http.StatusOK, payload)

	if err != nil {
		app.errorLog.Println(err)
	}
} 