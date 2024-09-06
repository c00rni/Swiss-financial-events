package main

import (
	"encoding/json"
	"net/http"
)

func respondWithoutContent(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, errorMessage string) error {
	type response struct {
		Error string `json:"error"`
	}

	return respondWithJSON(w, code, response{Error: errorMessage})
}

func decodeJSONBody(r *http.Request, ptr interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(ptr)
	if err != nil {
		return err
	}
	return nil
}
