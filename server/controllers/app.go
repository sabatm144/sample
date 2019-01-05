package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type appErr struct {
	Message string                 `json:"message"`
	Error   string                 `json:"error,omitempty"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

func renderJSON(w http.ResponseWriter, status int, res interface{}) {
	b, _ := json.MarshalIndent(res, "", "   ")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(b)
}

func renderERROR(w http.ResponseWriter, status int, message string) {
	renderJSON(w, status, appErr{Message: message})
}

func parseJSON(w http.ResponseWriter, params io.ReadCloser, data interface{}) bool {
	if params != nil {
		defer params.Close()
	}

	err := json.NewDecoder(params).Decode(data)
	log.Println(err)
	if err == nil {
		return true
	}

	renderERROR(w, http.StatusBadRequest, "Invalid JSON.")
	return false
}
