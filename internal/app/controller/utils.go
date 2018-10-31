package controller

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, r *http.Request, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		Error500(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
