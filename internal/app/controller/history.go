package controller

import (
	"net/http"
	"strconv"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/people"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/exclude_record"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
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
    t := checkin.CheckinID
    status := checkin.DefaultCheckiner.Status()

    if t == 0 {
	    respondJSON("", w, r)
        return 
    }

    record := checkin.GetCurrentPeopleSet()
    data, err := CheckinResult(&record, t)
    if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
    }

    data.Timestamp          = t 
	data.ExpectedCount, err = people.CountPeople()
    data.Status             = status
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
	} else if history == nil {
		Error(w, err, http.StatusNotFound)
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
    data.Status         = schema.STOPPED

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

    excludeSet := exclude_record.MakeExcludePeopleSet(exclude)

	return &schema.CheckinResp{
		ActualCount:   len(*record),
		Present:       dB2CheckinPeople(present, nil),
		Absent:        dB2CheckinPeople(absent, &excludeSet),
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

func dB2CheckinPeople(l []*schema.DBPerson, filter *map[string]int64) []*schema.CheckinPerson {
	res := []*schema.CheckinPerson{}
    if filter == nil {
	    res = make([]*schema.CheckinPerson, len(l))
        for i, p := range l {
            res[i] = p.CheckinPerson()
        }
    } else {
        for _, p := range l {
            if _, ok := (*filter)[p.NationalID]; !ok {
                res = append(res, p.CheckinPerson())
            }
        }
    }
	return res
}
