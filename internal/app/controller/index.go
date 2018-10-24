package controller

import (
	"io/ioutil"
	"net/http"
)

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("web/template/index.html")
	if err != nil {
		Error500(w, r)
		return
	}

	w.Write(body)
}
