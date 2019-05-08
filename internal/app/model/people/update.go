package people

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/images"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// AddPerson add a person to db
func AddPerson(p *schema.Person, imgs []string) error {
	if len(imgs) == 0 {
		return fmt.Errorf("should send at least one image")
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	personID := uuid.String()

	err = remote.Record(personID, imgs)
	if err != nil {
		return err
	}

	dbp := schema.NewDBPerson(p, personID)
	_, err = collection().InsertOne(context.Background(), dbp)
    if err != nil {
        return err
    }

    err = images.AddImages(personID, imgs)
    return err
}

// UpdatePerson update a person in db
func UpdatePerson(p *schema.DBPerson, imgs []string) error {
	if len(imgs) == 0 {
		return fmt.Errorf("should send at least one image")
	}

	err := remote.Record(p.ID, imgs)
	if err != nil {
		return err
	}

    updater := map[string]*schema.DBPerson{"$set": p}
	_, err = collection().UpdateOne(context.Background(), map[string]string{"_id": p.ID}, updater)
    if err != nil {
        return err
    }

    err = images.UpdateImages(p.ID, imgs)
    return err
}

// DeletePerson delete a person to db
func DeletePerson(id string) error {
	_, err := collection().DeleteOne(context.Background(), map[string]string{"_id": id})
	return err
}
