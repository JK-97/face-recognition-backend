package schema

import "gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"

// Person represents basic info of a person
type Person struct {
	ID           string `json:"id"`
	SerialNumber string `json:"serial_number"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	NationalID   string `json:"national_id"`
}

// DBPerson represents a person in db
type DBPerson struct {
	ID             string `json:"id" bson:"_id"`
	SerialNumber   string `json:"serial_number"  bson:"serial_number"`
	Name           string `json:"name" bson:"name"`
	Location       string `json:"location" bson:"location"`
	NationalID     string `json:"national_id" bson:"national_id"`
	CreateTime     int64  `json:"created_time" bson:"created_time"`
	LastUpdateTime int64  `json:"last_update_time" bson:"last_update_time"`
	Image          string `json:"image" bson:"image"`
}

// NewDBPerson creates DBPerson with Person
func NewDBPerson(p *Person, image string) *DBPerson {
	return &DBPerson{
		ID:             p.ID,
		SerialNumber:   p.SerialNumber,
		Name:           p.Name,
		Location:       p.Location,
		NationalID:     p.NationalID,
		CreateTime:     util.NowMilli(),
		LastUpdateTime: util.NowMilli(),
		Image:          image,
	}
}

// Person gets person info from DBPerson
func (p *DBPerson) Person() Person {
	return Person{
		ID:           p.ID,
		SerialNumber: p.SerialNumber,
		Name:         p.Name,
		Location:     p.Location,
		NationalID:   p.NationalID,
	}
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
