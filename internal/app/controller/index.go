package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func serveStatic(path string, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		Error500(w, r)
		return
	}

	w.Write(body)
}

// PingGet returns pong
func PingGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	serveStatic("web/template/index.html", w, r)
}

// ImgGET send imgs
func ImgGET(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/img/"):]
	serveStatic("img/"+title, w, r)
}
