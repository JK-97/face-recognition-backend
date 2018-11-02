package model

import (
	"fmt"
	"path/filepath"

	"github.com/satori/go.uuid"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

var people = map[string]*schema.DBPerson{}

// CountPeople counts people in db
func CountPeople() int {
	return len(people)
}

// GetPerson get a person by id in db
func GetPerson(id string) *schema.DBPerson {
	return people[id]
}

// GetPeople gets list of people in db
func GetPeople() []*schema.DBPerson {
	v := make([]*schema.DBPerson, 0, len(people))
	for _, person := range people {
		v = append(v, person)
	}
	return v
}

// AddPerson add a person to db
func AddPerson(p schema.Person, images []string) error {
	if len(images) == 0 {
		return fmt.Errorf("should send at least one image")
	}

	uuid, err := uuid.NewV1()
	if err != nil {
		return err
	}
	personID := uuid.String()

	err = remote.Record(personID, images)
	if err != nil {
		return err
	}

	img, err := base64Str2Img(images[0])
	fname := saveImg(img, "profile")

	p.ID = personID
	dbp := &schema.DBPerson{
		Person:         p,
		CreateTime:     util.NowMilli(),
		LastUpdateTime: util.NowMilli(),
		ImageURL:       "/v1/img/" + filepath.Base(fname),
	}
	people[p.ID] = dbp

	return nil
}
