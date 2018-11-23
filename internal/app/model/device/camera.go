package device

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

func collection() *mongo.Collection {
	return model.DB.Collection("device")
}

// GetCameras get cameralist in db
func GetCameras() ([]*schema.Camera, error) {

	ctx := context.Background()
	cur, err := collection().Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	result := []*schema.Camera{}
	for cur.Next(ctx) {
		dbp := &schema.Camera{}
		if err := cur.Decode(dbp); err != nil {
			return nil, err
		}
		result = append(result, dbp)
	}
	return result, nil
}

// AddCamera create camera in db
func AddCamera(h *schema.Camera) error {
	_, err := collection().InsertOne(context.Background(), h)
	return err
}
