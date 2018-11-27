package controller

import (
	"net/http"
	"strconv"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/people"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/exclude_record"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
)

// CheckinHistoryGET get checkin history timestamps
func CheckinHistoryGET(w http.ResponseWriter, r *http.Request) {
	ts, err := checkin.HistoryTimestamps(0, 0)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}
	respondJSON(schema.CheckinHistoryResp(ts), w, r)
}

// CheckinGETCurrent returns result in memory
func CheckinGETCurrent(w http.ResponseWriter, r *http.Request) {
    t := util.NowMilli()
    record := checkin.GetCurrentPeopleSet(false)
    data, err := CheckinResult(&record, t)
    if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
    }

    // data.Timestamp      = history.StartTime
    // data.CostTime       = history.EndTime - history.StartTime
	data.ExpectedCount, err = people.CountPeople()
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

	respondJSON(data, w, r)
}


// CheckinGETHistory returns result in db
func CheckinGETHistory(w http.ResponseWriter, r *http.Request, t int64) {
	history, err := checkin.GetHistory(t)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	if history == nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

    data, err := CheckinResult(&history.Record, t)
    if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
    }
    data.Timestamp      = history.StartTime
    data.CostTime       = history.EndTime - history.StartTime
    data.ExpectedCount  = history.ExpectedCount

	respondJSON(data, w, r)
}

// CheckinResult return checkin result under query condition
func CheckinResult(record *[]string, t int64) (*schema.CheckinResp, error) {
	present, err := people.GetPeople(people.NewFilterPresent(*record), 0, 0)
	if err != nil {
		return nil, err
	}

	absent, err := people.GetPeople(people.NewFilterAbsent(*record, t), 0, 0)
	if err != nil {
		return nil, err
	}

    // exclude_record
    exclude, err := exclude_record.GetExcludeRecord(exclude_record.NewFilterExcludeHistory(t), -1, -1)
	if err != nil {
		return nil, err
	}

	return &schema.CheckinResp{
		ActualCount:   len(*record),
		Present:       dB2CheckinPeople(present),
		Absent:        dB2CheckinPeople(absent),
        ExcludeRecord: exclude,
	}, nil
}

// CheckinGET get a checkin record
func CheckinGET(w http.ResponseWriter, r *http.Request) {
	tString := r.URL.Query().Get("timestamp")
    if tString == "" {
		CheckinGETCurrent(w, r)
		return
    }

	t, err := strconv.Atoi(tString)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}
    CheckinGETHistory(w, r, int64(t))
}

func dB2CheckinPeople(l []*schema.DBPerson) []*schema.CheckinPerson {
	res := make([]*schema.CheckinPerson, len(l))
	for i, p := range l {
		res[i] = p.CheckinPerson()
	}
	return res
}
