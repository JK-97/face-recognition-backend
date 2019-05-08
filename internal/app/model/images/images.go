package images

import (
	"fmt"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
)

func getImagePath(pid string) string {
    cfg := config.Config()
    appid := cfg.GetString("appid")
    prefix := cfg.GetString("app-name")
    if appid != "" {
        prefix = appid
    }
    return fmt.Sprintf("application/%s/database/people/%s", prefix, pid)
}

// GetImageIDs return a person image id list
func GetImageIDs(pid string)([]string, error) {
    appImgPath := getImagePath(pid)
    return util.ListImgs(appImgPath)
}

// GetImages return a person full images or single image
func GetImages(pid string, image_id string) ([]string, error) {
	var result []string

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
        result = append(result, b64Img)
    }
    return result, nil
}

// UpdateImages update a person's images in db
func UpdateImages(ID string, imgs []string) error {
    appImgPath := getImagePath(ID)
    for index, img := range imgs {
        img, err := util.Base64Str2Img(img)
        if err != nil {
            return err
        }
        util.SaveImg(img, appImgPath, index)
    }
    return nil
}

// AddImages add a person's image in db
func AddImages(ID string, imgs []string) error {
    return UpdateImages(ID, imgs)
}
