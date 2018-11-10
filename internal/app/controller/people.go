package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/people"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// FaceRecordsGET returns one captured image
func FaceRecordsGET(w http.ResponseWriter, r *http.Request) {
	b, err := remote.Capture()
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	face := schema.FaceRecordsGETResp{
		Img: b,
	}
	respondJSON(face, w, r)
}

// CheckinPeoplePOSTDELETE handles post and delete
// POST adds a person to checkin people list
// DELETE ?id=xxx delete a checkin people
func CheckinPeoplePOSTDELETE(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CheckinPeoplePOST(w, r)
	case http.MethodDelete:
		CheckinPeopleDELETE(w, r)
	default:
		http.NotFound(w, r)
	}
}

// CheckinPeoplePOST adds a person to checkin people list
func CheckinPeoplePOST(w http.ResponseWriter, r *http.Request) {
	if checkin.DefaultCheckiner.Status() == schema.CHECKING {
		Error(w, fmt.Errorf("Cannot add person while checking in"), http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	var p schema.CheckinPeoplePOSTReq
	err = json.Unmarshal(b, &p)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return
	}

	err = people.AddPerson(&p.Person, p.Images)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
}

// CheckinPeopleDELETE ?id=xxx delete a checkin people
func CheckinPeopleDELETE(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := people.DeletePerson(id)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
}

// CheckinPeopleListGET returns checkin people list
func CheckinPeopleListGET(w http.ResponseWriter, r *http.Request) {
	people, err := people.GetPeople(nil, 0, 0)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	respondJSON(people, w, r)
}

// StartRecordingPOST returns ok if ready to capture images
func StartRecordingPOST(w http.ResponseWriter, r *http.Request) {
}
