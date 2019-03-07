package images

import (
	"fmt"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/google/uuid"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

func getImagePath(pid string) string {
    cfg := config.Config()
    appid := cfg.GetString("appid")
    prefix := "local_test"
    if appid != "" {
        prefix = appid
    }
    return fmt.Sprintf("application/%s/database/people/%s", prefix, pid)
}

func collection() *mongo.Collection {
	return model.DB.Collection("images")
}

// GetImageIDs return a person image id list
func GetImageIDs(pid string)([]string, error) {
    return getImageIDsByCeph(pid)
}

func getImageIDsByCeph(pid string)([]string, error) {
    appImgPath := getImagePath(pid)
    return util.ListImgs(appImgPath)
}

func getImageIDsByMongo(pid string)([]string, error) {
    ctx := context.Background()
	cur, err := collection().Find(ctx, map[string]string{"pid": pid})
    if err != nil {
        return nil, err
    }

    var ids []string
	for cur.Next(ctx) {
		dbp := &schema.DBImagesOnlyID{}
		if err := cur.Decode(dbp); err != nil {
			return nil, err
		}
        ids = append(ids, dbp.ID)
    }
    return ids, nil
}

// GetImages return a person full images or single image
func GetImages(pid string, image_id string) (*schema.DBImages, error) {
    return getImagesByCeph(pid, image_id)
}

func getImagesByCeph(pid string, image_id string) (*schema.DBImages, error) {
	result := &schema.DBImages{
        PID: pid,
    }

    appImgPath := getImagePath(pid)
    var filenames []string

    if image_id == "" {
        tmpfiles, err := util.ListImgs(appImgPath)
        if err != nil {
            return nil, err
        }
        filenames = tmpfiles
    } else {
        filenames = append(filenames, image_id)
    }

    for _, f := range filenames {
        img, err := util.GetImg(fmt.Sprintf("%s/%s", appImgPath, f))
        if err != nil {
            return nil, err
        }
        b64Img, err := util.Img2Base64Str(img)
        if err != nil {
            return nil , err
        }
        result.Images = append(result.Images, b64Img)
    }
    return result, nil
}

func getFullImagesFilter(id string) map[string]string {
    return map[string]string{"pid": id}
}

func getSingleImageFilter(id string, imageID string) map[string]string {
    return map[string]string{
        "pid": id,
        "_id": imageID,
    }
}

func getImagesByMongo(id string, image_id string) (*schema.DBImages, error){
    var options map[string]string
    if image_id == "" {
        options = getFullImagesFilter(id)
    } else {
        options = getSingleImageFilter(id, image_id)
    }

    ctx := context.Background()
	cur, err := collection().Find(ctx, options)
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
    return setImagesByCeph(ID, imgs, true)
}

func setImagesByCeph(ID string, imgs []string, update bool) error {
    // upedataImg 更新图像, 优先删除remote， 然后修改本地cache， 异步上传
    appImgPath := getImagePath(ID)

    if update {
        err := util.RemoveImg(appImgPath, "")
        if err != nil {
            return err
        }
    }

    for index, img := range imgs {
        img, err := util.Base64Str2Img(img)
        if err != nil {
            return err
        }
        util.SaveImg(img, appImgPath, index)
    }
    return nil
}

func updateImagesByMongo(ID string, imgs []string) error {
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
    return setImagesByCeph(id, imgs, false)
}

func addImagesByMongo(id string, imgs []string) error {
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
