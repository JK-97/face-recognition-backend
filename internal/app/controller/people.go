package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/base64"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/checkin"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/exclude_record"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/people"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/images"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// FaceRecordsGET returns one captured image
func FaceRecordsGET(w http.ResponseWriter, r *http.Request) {
	deviceName := r.URL.Query().Get("device_name")
	b, err := remote.Capture(deviceName)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}
	face := schema.FaceRecordsGETResp{
		Img: b,
	}
	respondJSON(face, w, r)
}

// CheckinPeoplePOSTDELETEGETPUT handles post and delete
// POST adds a person to checkin people list
// DELETE ?id=xxx delete a checkin people
func CheckinPeoplePOSTDELETEGETPUT(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		CheckinPeoplePUT(w, r)
	case http.MethodGet:
		CheckinPeopleGET(w, r)
	case http.MethodPost:
		CheckinPeoplePOST(w, r)
	case http.MethodDelete:
		CheckinPeopleDELETE(w, r)
	default:
		http.NotFound(w, r)
	}
}

// CheckinPeopleGET returns person in db 
func CheckinPeopleGET(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
    if id == "" {
		Error(w, fmt.Errorf("Required args: id"), http.StatusBadRequest)
		return
    }

    dbPerson, err := people.GetPerson(id)
    if err != nil {
		Error(w, err, http.StatusNotFound)
		return
    }

	fullImage := r.URL.Query().Get("full_images")
	onlyImageID := r.URL.Query().Get("only_image_id")

    resp := schema.CheckinPeoplePOSTReq{
        Person: dbPerson.Person(),
    }

    if fullImage != "" {
        if onlyImageID == "" {
            imgs, err := images.GetImages(dbPerson.ID, images.GetFullImages(dbPerson.ID))
            if err != nil {
                Error(w, err, http.StatusNotFound)
                return
            }
            resp.Images = imgs.Images
        } else {
            ids, err := images.GetImageIDs(dbPerson.ID)
            if err != nil {
                Error(w, err, http.StatusNotFound)
                return
            }
            resp.ImageIDs = ids
        }
    }

	respondJSON(resp, w, r)
}

// CheckinPeoplePUT update person in db 
func CheckinPeoplePUT(w http.ResponseWriter, r *http.Request) {
	if checkin.DefaultCheckiner.Status() == schema.CHECKING {
		Error(w, fmt.Errorf("Cannot update person while checking in"), http.StatusBadRequest)
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

	err = people.UpdatePerson(&p.Person, p.Images)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

    err = images.UpdateImages(p.Person.ID, p.Images)
    if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
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

    err = images.AddImages(p.ID, p.Images)
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
	exclude := r.URL.Query().Get("exclude")

	people, err := people.GetPeople(nil, 0, 0)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
		return
	}

	var excludePeople = make(map[string]int64)
	if exclude != "" {
		excludePeople, err = exclude_record.GetExcludePeopleSetNow()
	}

	var peopleRet = []*schema.DBPerson{}
	for _, p := range people {
		if _, ok := excludePeople[p.ID]; !ok {
			peopleRet = append(peopleRet, p)
		}
	}

	respondJSON(peopleRet, w, r)
}

// StartRecordingPOST returns ok if ready to capture images
func StartRecordingPOST(w http.ResponseWriter, r *http.Request) {
}

// CheckinPeopleImageGET returns people image by id
func CheckinPeopleImageGET(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
    image_id := r.URL.Query().Get("image_id")

    var image string
    if image_id == "" {
        p, err := people.GetPerson(id)
        if err != nil {
            Error(w, err, http.StatusInternalServerError)
            return
        }
        image = p.Image
    } else {
        imgs, err := images.GetImages(id, images.GetSingleImage(id, image_id))
        if err != nil || len(imgs.Images) != 1 {
            Error(w, err, http.StatusInternalServerError)
            return
        }
        image = imgs.Images[0]
    }

    by, err := base64.StdEncoding.DecodeString(image)
    if err != nil {
        Error(w, err, http.StatusInternalServerError)
        return
    }
	w.Header().Set("Content-Type", "image/*")
	w.Write(by)
}
