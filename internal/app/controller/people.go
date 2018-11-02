package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// FaceRecordsGET returns one captured image
func FaceRecordsGET(w http.ResponseWriter, r *http.Request) {
	b, err := remote.Capture()
	if err != nil {
		http.Error(w, "cannot capture image", http.StatusInternalServerError)
		return
	}
	face := schema.FaceRecordsGETResp{
		Img: b,
	}
	respondJSON(face, w, r)
}

// CheckinPeoplePOST adds a person to checkin people list
func CheckinPeoplePOST(w http.ResponseWriter, r *http.Request) {
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

	err = model.AddPerson(p.Person, p.Images)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
}

// CheckinPeopleListGET returns checkin people list
func CheckinPeopleListGET(w http.ResponseWriter, r *http.Request) {
	people := model.GetPeople()
	respondJSON(people, w, r)
}

// StartRecordingPOST returns ok if ready to capture images
func StartRecordingPOST(w http.ResponseWriter, r *http.Request) {
}
