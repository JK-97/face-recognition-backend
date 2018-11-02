package schema

// Person represents base info of a person
type Person struct {
	ID           string `json:"id"`
	SerialNumber string `json:"serial_number"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	NationalID   string `json:"national_id"`
}

// DBPerson represents a person in db
type DBPerson struct {
	Person
	CreateTime     int64  `json:"created_time"`
	LastUpdateTime int64  `json:"last_update_time"`
	ImageURL       string `json:"image_url"`
}

// CheckinPeoplePOSTReq is request to CheckinPeoplePost
type CheckinPeoplePOSTReq struct {
	Person
	Images []string `json:"images"`
}

// FaceRecordsGETResp is response of FaceRecordsGET
type FaceRecordsGETResp struct {
	Img string `json:"img"`
}

// CheckinPeopleListResp is response of CheckinPeopleListGET
type CheckinPeopleListResp []DBPerson
