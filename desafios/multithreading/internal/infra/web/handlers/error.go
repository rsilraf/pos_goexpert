package handlers

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

// helper function
func sendError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	var outErr Error
	if err != nil {
		outErr = Error{err.Error()}
	}
	json.NewEncoder(w).Encode(outErr)
}
