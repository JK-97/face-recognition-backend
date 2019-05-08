package controller

import (
	"net/http"
	"strconv"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
)

func getURLInt64Params(param string, r *http.Request) (int64, error) {
    p := r.URL.Query().Get(param)
    ret, err := strconv.Atoi(p)
    return int64(ret), err
}

// CheckinGET get a checkin record
func CheckinGET(w http.ResponseWriter, r *http.Request) {
	t, err := getURLInt64Params("timestamp", r)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

    cameraID := r.URL.Query().Get("camera_id")
    records, err := checkin.GetHistoryRecords(t, t, cameraID)
    if err != nil {
		Error(w, err, http.StatusInternalServerError)
        return
    }
	respondJSON(records, w, r)
}

// CheckinRangeGET get a set of checkin record
func CheckinRangeGET(w http.ResponseWriter, r *http.Request) {
	startTime, err := getURLInt64Params("start_time", r)
    if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
    }

	endTime, err := getURLInt64Params("end_time", r)
    if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
    }

    cameraID := r.URL.Query().Get("camera_id")

    records, err := checkin.GetHistoryRecords(startTime, endTime, cameraID)
    if err != nil {
		Error(w, err, http.StatusInternalServerError)
        return
    }
	respondJSON(records, w, r)
}
