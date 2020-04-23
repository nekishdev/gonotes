package utils

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func Respond(w http.ResponseWriter, data interface{})  {
	w.Header().Add("Content-Type", "application/json")
	response := map[string]interface{}{
		"status": true,
		"data": data,
	}
	json.NewEncoder(w).Encode(response)
}

func Abort(w http.ResponseWriter, code int, data interface{})  {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	response := map[string]interface{}{
		"status": false,
		"data": data,
	}
	json.NewEncoder(w).Encode(response)
}

func GetUrlVar(r *http.Request, key string) string {
	params := mux.Vars(r)
	return params[key]
}