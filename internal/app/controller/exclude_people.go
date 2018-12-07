package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/exclude_record"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// ExcludeRecordGETPOSTPUT handle requests
func ExcludeRecordGETPOSTPUT(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ExcludeRecordPOST(w, r)
	case http.MethodGet:
		ExcludeRecordGET(w, r)
	case http.MethodPut:
		ExcludeRecordPUT(w, r)
	default:
		http.NotFound(w, r)
	}
}

// ExcludeRecordGET handle get requests
func ExcludeRecordGET(w http.ResponseWriter, r *http.Request) {
	excludeFlag := r.URL.Query().Get("all")
	skip := r.URL.Query().Get("skip")
	limit := r.URL.Query().Get("limit")
	s := -1
	l := -1
	if skip != "" {
		s, _ = strconv.Atoi(skip)
	}
	if limit != "" {
		l, _ = strconv.Atoi(limit)
	}

	filter := exclude_record.NewFilterExclude(util.NowMilli(), excludeFlag != "")
	resp, err := exclude_record.GetExcludeRecord(filter, l, s)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	respondJSON(resp, w, r)
}

// ExcludeRecordPOST handle get requests
func ExcludeRecordPOST(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	var excludeRecord schema.ExcludeRecordReq
	err = json.Unmarshal(b, &excludeRecord)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	people, err := exclude_record.GetExcludePeopleSetNow()
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

    log.Info("people: %v", people)
	for _, r := range excludeRecord.People {
		if _, ok := people[r.ID]; ok {
			Error(w, err, http.StatusNotAcceptable)
			return
		}
	}

	err = exclude_record.AddExcludeRecord(&excludeRecord, util.NowMilli())
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
}

// ExcludeRecordPUT handle delete requests
func ExcludeRecordPUT(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := exclude_record.UpdateRecord(id)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
}
