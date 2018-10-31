package controller

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
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

func serveStatic(path string, w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		Error500(w, r)
		return err
	}

	w.Write(body)
	return nil
}

func serveTemplate(path string, data interface{}, w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		Error500(w, r)
		return err
	}

	t, err := template.New(path).Parse(string(body))
	if err != nil {
		Error500(w, r)
		return err
	}

	err = t.Execute(w, data)
	if err != nil {
		Error500(w, r)
		return err
	}

	return nil
}
