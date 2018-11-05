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
	respondJSON(schema.CheckinHistoryResp(checkin.HistoryTimestamps()), w, r)
}

// CheckinGET get a checkin record
func CheckinGET(w http.ResponseWriter, r *http.Request) {
	tString := r.URL.Query().Get("timestamp")
	t, err := strconv.Atoi(tString)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	history := checkin.GetHistory(int64(t))
	if history == nil {
		http.Error(w, "requested checkin data is not available", http.StatusBadRequest)
		return
	}

	details := make([]schema.CheckinPerson, 0, len(history.Record))
	for id := range history.Record {
		person := schema.Person{ID: id}
		dbp := people.GetPerson(id)
		if dbp != nil {
			person = dbp.Person
		}
		details = append(details, schema.CheckinPerson{Person: person})
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
