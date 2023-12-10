package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type envelope map[string]any //envelops the json response under a key

// Credit: Alex Edwards, Let's Go Further
// Additional note. keeping the credit since I'm writing to a public github
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t") //encode the into json
	if err != nil {
		return err
	}

	js = append(js, '\n') //append a newline for neatness

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json") //set the header
	w.WriteHeader(status)                              //write status code
	w.Write(js)                                        //write the response

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		// Custom Error Handling: Alex Edwards, Let's Go Further Chapter 4
		//Additional note. keeping the credit since I'm writing to a public github
		return err
	}

	err := dec.Decode(&struct{}{}) //Decode the item into a struct
	if err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}
