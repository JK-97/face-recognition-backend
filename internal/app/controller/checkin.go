package controller

import (
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// CheckStatusGET returns checkin status "checking/stopped"
func CheckStatusGET(w http.ResponseWriter, r *http.Request) {
	status := checkin.Status()
	respondJSON(schema.CheckStatusResp{Status: status}, w, r)
}

// StartCheckinPOST starts check in
func StartCheckinPOST(w http.ResponseWriter, r *http.Request) {
	err := checkin.StartCheckin()
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}
}

// StopCheckinPOST stops check in
func StopCheckinPOST(w http.ResponseWriter, r *http.Request) {
	t, err := checkin.StopCheckin()
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}
	respondJSON(schema.StopCheckinResp{Timestamp: t}, w, r)
}
