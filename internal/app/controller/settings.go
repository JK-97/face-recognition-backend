package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/settings"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// SettingsPOSTGET handle settings about runtime
func SettingsPOSTGET(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		SettingsPOST(w, r)
	case http.MethodGet:
		SettingsGET(w, r)
	default:
		http.NotFound(w, r)
	}
}

// SettingsGET return camera list
func SettingsGET(w http.ResponseWriter, r *http.Request) {
	s, err := settings.GetSettings()
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	respondJSON(s, w, r)
}

// SettingsPOST add camera in db(test)
func SettingsPOST(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	var p schema.SettingsReq
	err = json.Unmarshal(b, &p)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}
	err = settings.UpdateSettings(&p)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
}
