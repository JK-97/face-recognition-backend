package images

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/google/uuid"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

func collection() *mongo.Collection {
	return model.DB.Collection("images")
}

// GetImages return a person full images
func GetImages(id string) (*schema.DBImages, error){
    ctx := context.Background()
	cur, err := collection().Find(ctx, map[string]string{"pid": id})
    if err != nil {
        return nil, err
    }
	defer cur.Close(ctx)

	result := &schema.DBImages{
        PID: id,
    }

	for cur.Next(ctx) {
		dbp := &schema.DBImages{}
		if err := cur.Decode(dbp); err != nil {
			return nil, err
		}
        result.Images = append(result.Images, dbp.Image)
    }
    return result, nil
}

// UpdateImages update a person's images in db
func UpdateImages(ID string, imgs []string) error {
    for index, img := range imgs {
        updater := map[string]schema.DBImagesUpdater{"$set": schema.DBImagesUpdater{
            PID:        ID,
            Image:      img,
        }}
        _, err := collection().UpdateOne(context.Background(), map[string]string{"pid": ID, "image_id": string(index)}, updater)
        if err != nil {
            return err
        }
    }
    return nil
}

// AddImages add a person's image in db
func AddImages(id string, imgs []string) error {
    for index, img := range imgs {
        uid, _ := uuid.NewUUID()
        doc := schema.DBImages{
            ID:             uid.String(),
            PID:            id,
            Image:          img,
            ImageID:        string(index),
        }
        _, err := collection().InsertOne(context.Background(), doc)
        if err != nil {
            return err
        }
    }
    return nil
}
