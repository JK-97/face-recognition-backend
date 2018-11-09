package controller

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// Error handles server error
func Error(w http.ResponseWriter, err error, code int) {
	if h, ok := err.(schema.HTTPError); ok {
		code = h.Code
	}

	http.Error(w, http.StatusText(code), code)
	log.Error(err)
}

func respondJSON(obj interface{}, w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(obj)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func serveStatic(path string, w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return err
	}

	w.Write(body)
	return nil
}

func serveTemplate(path string, data interface{}, w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return err
	}

	t, err := template.New(path).Parse(string(body))
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return err
	}

	err = t.Execute(w, data)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return err
	}

	return nil
}
