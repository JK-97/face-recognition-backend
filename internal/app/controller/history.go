package controller

import (
	"net/http"
	"strconv"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/people"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// CheckinHistoryGET get checkin history timestamps
func CheckinHistoryGET(w http.ResponseWriter, r *http.Request) {
	ts, err := checkin.HistoryTimestamps(10, 0)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}
	respondJSON(schema.CheckinHistoryResp(ts), w, r)
}

// CheckinGET get a checkin record
func CheckinGET(w http.ResponseWriter, r *http.Request) {
	tString := r.URL.Query().Get("timestamp")
	t, err := strconv.Atoi(tString)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	history, err := checkin.GetHistory(int64(t))
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}
	if history == nil {
		http.Error(w, "requested checkin data is not available", http.StatusBadRequest)
		return
	}

	details := make([]schema.CheckinPerson, 0, len(history.Record))
	for _, id := range history.Record {
		person := schema.CheckinPerson{ID: id, Name: "unknown"}
		dbp, err := people.GetPerson(id)
		if dbp != nil {
			person = dbp.CheckinPerson()
		}
		if err != nil {
			Error(w, err, http.StatusInternalServerError)
			return
		}
		details = append(details, person)
	}
	data := schema.CheckinResp{
		Timestamp:     history.StartTime,
		CostTime:      history.EndTime - history.StartTime,
		ExpectedCount: history.ExpectedCount,
		ActualCount:   history.ActualCount,
		Detail:        details,
	}
	respondJSON(data, w, r)
}
