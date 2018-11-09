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
		Error(w, err, http.StatusInternalServerError)
		return
	}
	if history == nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

	present, err := people.GetPeople(people.NewFilterPresent(history.Record), 10, 0)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
	}

	absent, err := people.GetPeople(people.NewFilterAbsent(history.Record, int64(t)), 10, 0)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
	}

	data := schema.CheckinResp{
		Timestamp:     history.StartTime,
		CostTime:      history.EndTime - history.StartTime,
		ExpectedCount: history.ExpectedCount,
		ActualCount:   history.ActualCount,
		Present:       dB2CheckinPeople(present),
		Absent:        dB2CheckinPeople(absent),
	}
	respondJSON(data, w, r)
}

func dB2CheckinPeople(l []*schema.DBPerson) []*schema.CheckinPerson {
	res := make([]*schema.CheckinPerson, len(l))
	for i, p := range l {
		res[i] = p.CheckinPerson()
	}
	return res
}
