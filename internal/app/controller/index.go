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
	serveStatic("web/template/index.html", w, r)
}

// ImgGET send imgs
func ImgGET(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/img/"):]
	serveStatic("img/"+title, w, r)
}
