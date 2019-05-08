package schema

// Person represents basic info of a person
type Person struct {
	SerialNumber string `json:"serial_number"`
}

// DBPerson represents a person in db
type DBPerson struct {
	ID             string `json:"id" bson:"_id"`
	SerialNumber   string `json:"serial_number"  bson:"serial_number"`
}

// NewDBPerson creates DBPerson with Person
func NewDBPerson(p *Person, id string) *DBPerson {
	return &DBPerson{
		ID:             id,
		SerialNumber:   p.SerialNumber,
	}
}

// Person gets person info from DBPerson
func (p *DBPerson) Person() *Person {
	return &Person{
		SerialNumber: p.SerialNumber,
	}
}

// CheckinPeoplePOSTReq is request to CheckinPeoplePost
type CheckinPeoplePOSTReq struct {
    Person      Person   `json:"person"`
	Images      []string `json:"images"`
}

// CheckinPeoplePUTReq is request to CheckinPeoplePUT
type CheckinPeoplePUTReq struct {
    Person      DBPerson `json:"person"`
	Images      []string `json:"images"`
}

// CheckinPeopleGETResp is response to CheckingPeopleGET
type CheckinPeopleGETResp struct {
    Person      DBPerson `json:"person"`
	Images      []string `json:"images"`
    ImageIDs    []string `json:"image_ids"`
}

// FaceRecordsGETResp is response of FaceRecordsGET
type FaceRecordsGETResp struct {
	Img string `json:"img"`
}

// CheckinPeopleListResp is response of CheckinPeopleListGET
type CheckinPeopleListResp []DBPerson
