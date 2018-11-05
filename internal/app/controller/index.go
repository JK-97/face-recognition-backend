package controller

import (
	"fmt"
	"net/http"
)

// PingGet returns pong
func PingGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveStatic("web/template/index.html", w, r)
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// ImgGET send imgs
func ImgGET(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/v1/img/"):]
	serveStatic("img/"+title, w, r)
}
