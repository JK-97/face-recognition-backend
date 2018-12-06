package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/device"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// CamerasGETPOST handle camera list
func CamerasGETPOST(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CamerasPOST(w, r)
	case http.MethodGet:
		CamerasGET(w, r)
	default:
		http.NotFound(w, r)
	}
}

// CamerasGET return camera list
func CamerasGET(w http.ResponseWriter, r *http.Request) {
	cameras, err := device.GetCameras()
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	respondJSON(cameras, w, r)
}

// CamerasPOST add camera in db(test)
func CamerasPOST(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	var p schema.Camera
	err = json.Unmarshal(b, &p)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}
	err = device.AddCamera(&p)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
}
