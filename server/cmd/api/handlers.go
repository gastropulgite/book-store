package main

import (
	"net/http"
)

type cridentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type jsonResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	
	var creds cridentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds) 

	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "Invalid JSON"
		app.writeJSON(w, http.StatusBadRequest, payload)
	}
	// err := json.NewDecoder(r.Body).Decode(&creds)

	// if err != nil {
	// 	app.errorLog.Println("Invalid json!")

	// 	payload.Error = true
	// 	payload.Message = "Invalid json!"
	// 	out, err := json.MarshalIndent(payload, "", "\t")
		
	// 	if err != nil {
	// 		app.errorLog.Println(err)
	// 	}
		
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write(out)
	// 	return
	// } 

	app.infoLog.Println(creds.UserName, creds.Password)

	payload.Error = false
	payload.Message = "Signed In!"

	err = app.writeJSON(w, http.StatusOK, payload)

	if err != nil {
		app.errorLog.Println(err)
	}
} 