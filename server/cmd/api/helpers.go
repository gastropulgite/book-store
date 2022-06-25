package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)


func (app *application) readJSON(
	w http.ResponseWriter, 
	r *http.Request, 
	data interface{}) error  {

	maxBytes := 1048576 // 1Mb
	r.Body =  http.MaxBytesReader(w, r.Body, int64(maxBytes))	//limit the body size from client
	
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}
	
	err = dec.Decode(&struct{}{}) // check for single json value
	if err != nil {
		return errors.New("Body must have only a single JSON value!")
	}

	return nil
}

func (app *application) writeJSON(
	w http.ResponseWriter, 
	status int, 
	data interface{}, 
	headers ...http.Header) error  {
	
	out, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {

			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	
	if err != nil {
		return err
	}
	return nil
}

func (app * application) errorJSON(w http.ResponseWriter, err error, status ...int) {

	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var customError error

	switch {
	case strings.Contains(err.Error(), `SQLSTATE 23505`):
		customError = errors.New(`Duplicate value violates constraint`)
		statusCode = http.StatusForbidden

	case strings.Contains(err.Error(), `SQLSTATE 22001`):
		customError = errors.New(`The value you are trying to insert to too large!`)
		statusCode = http.StatusForbidden

	case strings.Contains(err.Error(), `SQLSTATE 23503`):
		customError = errors.New(`Foreign key violation!`)
		statusCode = http.StatusForbidden

	default:
		customError = err
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = customError.Error()

	app.writeJSON(w, statusCode, payload)

	return
}